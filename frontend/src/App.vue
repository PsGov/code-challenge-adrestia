<template>
  <div id="app">
    <h1>User Management</h1>
    <UserList />
  </div>
</template>

<script>
import UserList from './components/UserList.vue';

export default {
  name: "App",
  components: {
    UserList,
  },
};
</script>
