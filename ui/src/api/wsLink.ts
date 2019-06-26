import { ApolloLink, Operation, FetchResult, Observable } from 'apollo-link';

import { SubscriptionClient, ClientOptions } from 'subscriptions-transport-ws/dist/client';

export class WebSocketLink extends ApolloLink {
  private subscriptionClient: SubscriptionClient;

  constructor(uri: string, options?: ClientOptions) {
    super();
    this.subscriptionClient = new SubscriptionClient(uri, options, WebSocket);
  }

  public request(operation: Operation): Observable<FetchResult> | null {
    return this.subscriptionClient.request(operation) as Observable<FetchResult>;
  }
}
