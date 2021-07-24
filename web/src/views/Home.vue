<template lang="pug">
van-loading(v-if="currentUserLoading || listMessagesLoading ", style="text-align: center; margin-top: 10px") Loading...
van-cell(v-else, v-for="(msg, i) in messages", :key="msg.id", :title="msg.text")
van-field(v-model="createMessageState")
  template(#button)
    van-button(size="small", :loading="createMessageLoading" @click="createMessage") add
</template>

<script>
import { ref } from "vue";
import { useQuery, useResult, useMutation } from "@vue/apollo-composable";
import gql from "graphql-tag";

const CURRENT_USER = gql`
  query currentUser {
    currentUser {
      id
      name
      avatarUrl
    }
  }
`;

const LIST_MESSAGES = gql`
  query listMessages {
    messages {
      id
      text
      createdAt
      user {
        id
        name
        avatarUrl
      }
    }
  }
`;

const CREATE_MESSAGE = gql`
  mutation createMessage($text: String!) {
    createMessage(input: { text: $text }) {
      id
      text
      createdAt
      user {
        id
        name
        avatarUrl
      }
    }
  }
`;

const ON_MESSAGECREATED = gql`
  subscription onMessageCreated {
    messageCreated {
      id
      text
      createdAt
      user {
        id
        name
        avatarUrl
      }
    }
  }
`;

export default {
  name: "Home",
  setup() {
    const { loading: currentUserLoading } = useQuery(CURRENT_USER);
    const {
      result: listMessages,
      loading: listMessagesLoading,
      subscribeToMore: onMessageCreated,
    } = useQuery(LIST_MESSAGES);
    const messages = useResult(listMessages, [], (data) => data.messages);
    const createMessageState = ref("");

    const {
      mutate: createMessage,
      loading: createMessageLoading,
      onDone: onCreateMessageDone,
    } = useMutation(CREATE_MESSAGE, () => ({
      variables: {
        text: createMessageState.value,
      },
    }));
    onCreateMessageDone(() => (createMessageState.value = ""));

    onMessageCreated(() => ({
      document: ON_MESSAGECREATED,
      updateQuery: (previousResult, { subscriptionData }) => {
        console.log(subscriptionData.data.messageCreated);
        previousResult.messages.push(subscriptionData.data.messageCreated);
        return previousResult;
      },
    }));

    return {
      currentUserLoading,
      listMessagesLoading,
      messages,
      createMessageState,
      createMessage,
      createMessageLoading,
    };
  },
};
</script>
