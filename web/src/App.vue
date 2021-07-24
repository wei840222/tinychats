<template lang="pug">
router-view
</template>

<script>
import { provide } from "vue";
import {
  ApolloClient,
  ApolloLink,
  createHttpLink,
  InMemoryCache,
  from,
  split,
} from "@apollo/client/core";
import { WebSocketLink } from "@apollo/client/link/ws";
import { getMainDefinition } from "@apollo/client/utilities";
import { createPersistedQueryLink } from "@apollo/client/link/persisted-queries";
import { sha256 } from "crypto-hash";
import { DefaultApolloClient } from "@vue/apollo-composable";
import liff from "@line/liff";

export default {
  setup() {
    if (process.env.NODE_ENV === "production") {
      liff.init({ liffId: "1656247924-eX5ZOvN0" }).then(() => {
        if (!liff.isLoggedIn()) {
          liff.login();
        }
      });
    }

    const authMiddleware = new ApolloLink((operation, forward) => {
      operation.setContext(({ headers = {} }) => ({
        headers: {
          ...headers,
          Authorization:
            process.env.NODE_ENV !== "production"
              ? process.env.VUE_APP_ACCESS_TOKEN
              : liff.isLoggedIn()
              ? liff.getAccessToken()
              : null,
        },
      }));
      return forward(operation);
    });

    const apolloClient = new ApolloClient({
      link: split(
        ({ query }) => {
          const definition = getMainDefinition(query);
          return (
            definition.kind === "OperationDefinition" &&
            definition.operation === "subscription"
          );
        },
        new WebSocketLink({
          uri: `wss://tinychats.herokuapp.com/graphql`,
          options: {
            reconnect: true,
          },
        }),
        from([
          authMiddleware,
          createPersistedQueryLink({
            useGETForHashedQueries: true,
            sha256,
          }),
          createHttpLink({
            uri: "https://tinychats.herokuapp.com/graphql",
          }),
        ])
      ),
      cache: new InMemoryCache(),
    });

    provide(DefaultApolloClient, apolloClient);
  },
};
</script>

<style lang="sass">
#app
  font-family: Avenir, Helvetica, Arial, sans-serif
  -webkit-font-smoothing: antialiased
  -moz-osx-font-smoothing: grayscale
  color: #2c3e50
</style>
