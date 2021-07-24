<template lang="pug">
van-loading(v-if="currentUserLoading || listMessagesLoading ", style="text-align: center; margin-top: 10px") Loading...
van-cell(v-else, v-for="(msg, i) in messages", :key="msg.id", :title="msg.text")
van-cell(v-for="(msg, i) in messagesCreated", :key="msg.id", :title="msg.text")
#message-end
van-field.fixedbutton(v-model="createMessageState")
  template(#button)
    van-button(size="small", :loading="createMessageLoading" @click="createMessage") send
</template>

<script>
import { ref, watch } from "vue";
import {
  useQuery,
  useResult,
  useMutation,
  useSubscription,
} from "@vue/apollo-composable";
import gql from "graphql-tag";
import jump from "jump.js";

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
      onResult: onListMessages,
    } = useQuery(LIST_MESSAGES);
    const messages = useResult(listMessages, [], (data) => data.messages);
    const messagesCreated = ref([]);
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

    const { result: onMessageCreated } = useSubscription(ON_MESSAGECREATED);

    watch(
      onMessageCreated,
      (data) => {
        messagesCreated.value.push(
          JSON.parse(JSON.stringify(data.messageCreated))
        );
        jump("#message-end");
      },
      {
        lazy: true,
      }
    );

    onListMessages(() => jump(window.innerHeight));

    return {
      currentUserLoading,
      listMessagesLoading,
      messages,
      messagesCreated,
      createMessageState,
      createMessage,
      createMessageLoading,
    };
  },
};
</script>

<style lang="sass" scoped>
.fixedbutton
    position: fixed
    bottom: 0px
    right: 0px
</style>
