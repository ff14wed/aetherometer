module.exports = {
  "devServer": {
    "port": 9080
  },
  "productionSourceMap": false,
  "pages": {
    "index": {
      "entry": "src/main.ts",
      "template": "public/index.html",
      "filename": "index.html"
    },
    "loading": {
      "entry": "src/loading.ts",
      "template": "public/index.html",
      "filename": "loading.html"
    }
  },
  "pluginOptions": {
    "electronBuilder": {
      "builderOptions": {
        "win": {
          "target": "zip"
        },
        "directories": {
          "output": "../dist"
        },
        "extraFiles": [
          {
            "from": "../resources/${os}",
            "to": "resources/bin",
            "filter": [
              "**/*"
            ]
          },
          {
            "from": "../resources/datasheets",
            "to": "resources/datasheets",
            "filter": [
              "**/*.csv"
            ]
          }
        ]
      }
    }
  },
  "transpileDependencies": [
    "vuetify"
  ]
}