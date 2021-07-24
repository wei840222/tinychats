<template lang="pug">
van-loading(v-if="loading", style="text-align: center; margin-top: 10px") Loading...
van-cell(v-else, v-for="(t, i) in todos", :key="t.id", :title="t.text")
  template(#right-icon)
    van-switch(v-model="doneSwitchState[t.id]", size="24")
van-field(v-model="createTodoState")
  template(#button)
    van-button(size="small", :loading="createTodoLoading" @click="createTodo") add
</template>

<script>
import { ref } from "vue";
import { useQuery, useResult, useMutation } from "@vue/apollo-composable";
import gql from "graphql-tag";

const LIST_TODOS = gql`
  query listTodos {
    todos {
      id
      text
      done
    }
  }
`;

export default {
  name: "Home",
  setup() {
    const { result, loading } = useQuery(LIST_TODOS);
    const todos = useResult(result, [], (data) => data.todos);
    const doneSwitchState = ref({});
    const createTodoState = ref("");

    const {
      mutate: createTodo,
      loading: createTodoLoading,
      onDone: onCreateTodoDone,
    } = useMutation(
      gql`
        mutation createTodo($text: String!) {
          createTodo(input: { text: $text }) {
            id
            text
            done
          }
        }
      `,
      () => ({
        variables: {
          text: createTodoState.value,
        },
        update: (cache, { data: { createTodo } }) => {
          let data = cache.readQuery({ query: LIST_TODOS });
          data = JSON.parse(JSON.stringify(data));
          data.todos.push(createTodo);
          cache.writeQuery({ query: LIST_TODOS, data });
        },
      })
    );
    onCreateTodoDone(() => (createTodoState.value = ""));
    return {
      loading,
      todos,
      doneSwitchState,
      createTodoState,
      createTodo,
      createTodoLoading,
    };
  },
};
</script>
