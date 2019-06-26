import * as gql from './gql';

import { ApolloLink, execute, makePromise, GraphQLRequest } from 'apollo-link';
import { HttpLink } from 'apollo-link-http';
import { fetch } from 'apollo-env';
import { ExecutionResultDataDefault } from 'graphql/execution/execute';

import { WebSocketLink } from './wsLink';
import { onError } from 'apollo-link-error';

const errHandlerLink = onError(({ graphQLErrors, networkError }) => {
  if (graphQLErrors) {
    graphQLErrors.map((e) => {
      throw e;
    });
  }
  if (networkError) { throw networkError; }
});

const executeOperation = (
  httpLink: ApolloLink, operation: GraphQLRequest, credentials: string,
) => {
  operation.context = {
    headers: {
      authorization: `Bearer ${credentials}`,
    },
  };

  return new Promise<{ [key: string]: any }>((resolve, reject) => {
    makePromise(execute(httpLink, operation))
      .then((result: ExecutionResultDataDefault) => {
        if (result.data) {
          resolve(result.data);
        } else {
          reject(result.errors);
        }
      })
      .catch((error) => reject(error));
  });
};


export default class GQLClient {
  private httpLink: ApolloLink;
  private wsLink: ApolloLink;
  private adminToken = '';

  constructor(port: number) {
    this.httpLink = ApolloLink.from([
      errHandlerLink,
      new HttpLink({ uri: `http://localhost:${port}/query`, fetch }),
    ]);

    this.wsLink = ApolloLink.from([
      errHandlerLink,
      new WebSocketLink(`ws://localhost:${port}/query`, {
        lazy: true,
        reconnect: true,
        connectionParams: () => ({
          authorization: this.adminToken,
        }),
      }),
    ]);
  }

  public initializeAdminToken = async (token: string, isAdminToken?: boolean): Promise<string> => {
    if (isAdminToken) {
      this.adminToken = token;
      return token;
    }
    const data = await executeOperation(
      this.httpLink, { query: gql.createAdminTokenMutation }, token,
    );
    this.adminToken = data.createAdminToken;
    return this.adminToken;
  }

  public subscribeToStreamEvents = (
    handleStreamEvent: (streamID: number, type: string) => void,
  ) => {
    return execute(this.wsLink, { query: gql.streamSubscription })
      .subscribe({
        next: (subscriptionData) => {
          if (!subscriptionData.data) { return; }
          const streamEvent = subscriptionData.data.streamEvent;
          const streamID = streamEvent.streamID;
          const { __typename } = streamEvent.type;
          handleStreamEvent(streamID, __typename);
        },
      });
  }

  public listStreams = async (): Promise<Array<{ id: number }>> => {
    const data = await executeOperation(
      this.httpLink, { query: gql.listStreamsQuery }, this.adminToken,
    );
    return data.streams;
  }

  public addPlugin = async (pluginURL: string): Promise<string> => {
    const variables = { pluginURL };

    const data = await executeOperation(
      this.httpLink, { query: gql.addPluginMutation, variables }, this.adminToken,
    );
    return data.addPlugin;
  }

  public removePlugin = async (apiToken: string): Promise<boolean> => {
    const variables = { apiToken };

    const data = await executeOperation(
      this.httpLink, { query: gql.removePluginMutation, variables }, this.adminToken,
    );
    return data.removePlugin;
  }
}
