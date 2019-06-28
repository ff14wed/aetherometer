import gql from 'graphql-tag';

export const streamSubscription = gql`
  subscription Streams {
    streamEvent {
      streamID
      type {
        __typename
      }
    }
  }
`;

export const listStreamsQuery = gql`
  query AllStreams {
    streams {
      id
    }
  }
`;

export const versionQuery = gql`
  query Version {
    apiVersion
  }
`;

export const createAdminTokenMutation = gql`
  mutation CreateAdminToken {
    createAdminToken
  }
`;

export const addPluginMutation = gql`
  mutation AddPlugin($pluginURL: String!) {
    addPlugin(pluginURL: $pluginURL)
  }
`;

export const removePluginMutation = gql`
  mutation RemovePlugin($apiToken: String!) {
    removePlugin(apiToken: $apiToken)
  }
`;
