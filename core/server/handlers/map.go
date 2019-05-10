package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/ff14wed/sibyl/backend/config"
	"go.uber.org/zap"
)

// MapHandler serves requests for map images. It will attempt to serve them
// from disk first, downloading from xivapi if necessary. Users are subject to
// Client IP based rate limiting.
//
//   Usage:
//   GET /maps/{map_id} -> Returns an image with content type image/png or
//   image/jpg
//
//   Example:
//   GET /maps/123 -> 200 OK, image/png
//   GET /maps/1234 -> 404 Not Found
type MapHandler struct {
	prefix string
	c      config.MapConfig
	logger *zap.Logger

	etagCache map[string]struct{}
}

// NewMapHandler creates a new MapHandler
//
// prefix is the path to this handler.
func NewMapHandler(prefix string, c config.Config, l *zap.Logger) *MapHandler {
	mapConfig := c.Maps
	if mapConfig.APIPath == "" {
		mapConfig.APIPath = "https://xivapi.com"
	}
	return &MapHandler{
		prefix:    prefix,
		c:         mapConfig,
		logger:    l.Named("map-handler"),
		etagCache: make(map[string]struct{}),
	}
}

// ServeHTTP serves maps for the handler
func (h *MapHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	escapedPath := r.URL.EscapedPath()
	sMapID := strings.TrimPrefix(escapedPath, h.prefix)
	if len(sMapID) == len(escapedPath) {
		h.logger.Debug("Error handling request path", zap.String("path", escapedPath))
		http.NotFound(w, r)
		return
	}
	mapID, err := strconv.ParseUint(sMapID, 10, 16)
	if err != nil {
		h.logger.Debug("Error parsing map ID", zap.String("path", escapedPath))
		http.NotFound(w, r)
		return
	}

	etag := fmt.Sprintf(`W/"map-%s"`, sMapID) // weak ETag validator

	if match := r.Header.Get("If-None-Match"); match != "" {
		if _, found := h.etagCache[etag]; found && strings.Contains(match, etag) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	contentType, imageBytes, cacheErr := getMapFromCache(h.c.Cache, mapID)
	if cacheErr != nil {
		var apiErr error
		contentType, imageBytes, apiErr = getMapFromAPI(h.c.APIPath, mapID)
		if apiErr != nil {
			h.logger.Debug("Unable to find map in cache or API",
				zap.NamedError("cacheErr", cacheErr),
				zap.NamedError("apiErr", apiErr),
			)
			http.NotFound(w, r)
			return
		}
		saveFilePath := path.Join(h.c.Cache, strconv.FormatUint(mapID, 10)+".jpg")
		saveErr := ioutil.WriteFile(saveFilePath, imageBytes, 0777)
		if saveErr != nil {
			h.logger.Error("Unable to cache file",
				zap.String("path", saveFilePath),
				zap.Error(saveErr),
			)
		}
	}
	w.Header().Set("Content-Type", contentType)

	w.Header().Set("ETag", etag)
	w.Header().Set("Cache-Control", "max-age=86400") // 1 day
	h.etagCache[etag] = struct{}{}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(imageBytes)
}

// getMapFromCache attempts to get the file from cache. It returns an error if
// the file is not found.
func getMapFromCache(cachePath string, mapID uint64) (string, []byte, error) {
	pathNoExt := path.Join(cachePath, strconv.FormatUint(mapID, 10))
	pngFile := pathNoExt + ".png"
	jpgFile := pathNoExt + ".jpg"

	pngExists, err := exists(pngFile)
	if err != nil {
		return "", nil, err
	}
	jpgExists, err := exists(jpgFile)
	if err != nil {
		return "", nil, err
	}
	// Content type based on extension... Don't need to be too fancy with
	// file type checking because I don't expect
	var (
		resolvedPath string
		contentType  string
	)
	if pngExists {
		contentType = "image/png"
		resolvedPath = pngFile
	} else if jpgExists {
		contentType = "image/jpeg"
		resolvedPath = jpgFile
	}
	if len(resolvedPath) == 0 {
		return "", nil, errors.New("not found")
	}

	imageBytes, err := ioutil.ReadFile(resolvedPath)
	if err != nil {
		return "", nil, err
	}

	return contentType, imageBytes, nil
}

// exists returns whether or not the path exists.
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

// getMapFromAPI queries the API and downloads the map
func getMapFromAPI(apiPath string, mapID uint64) (string, []byte, error) {
	resp, err := http.Get(fmt.Sprintf("%s/Map/%d", apiPath, mapID))
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, errors.New("not found")
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}

	mapURL, err := resolveMapURL(apiPath, bodyBytes)
	if err != nil {
		return "", nil, err
	}

	mapResp, err := http.Get(mapURL.String())
	if err != nil {
		return "", nil, err
	}
	defer mapResp.Body.Close()

	if mapResp.StatusCode != http.StatusOK {
		return "", nil, errors.New("not found")
	}

	contentType := mapResp.Header.Get("Content-Type")

	mapBytes, err := ioutil.ReadAll(mapResp.Body)
	if err != nil {
		return "", nil, err
	}

	return contentType, mapBytes, nil
}

func resolveMapURL(apiPath string, apiJSON []byte) (*url.URL, error) {
	data := make(map[string]json.RawMessage)

	err := json.Unmarshal(apiJSON, &data)
	if err != nil {
		return nil, err
	}

	filenameJSON, found := data["MapFilename"]
	if !found {
		return nil, fmt.Errorf("invalid response from API: %s", string(apiJSON))
	}
	var mapFilename string
	err = json.Unmarshal([]byte(filenameJSON), &mapFilename)
	if err != nil {
		return nil, err
	}

	mapFilename = strings.Replace(mapFilename, "\\", "", -1)
	mapPath, err := url.Parse(mapFilename)
	if err != nil {
		return nil, err
	}

	apiURL, err := url.Parse(apiPath)
	if err != nil {
		return nil, err
	}
	return apiURL.ResolveReference(mapPath), nil
}
