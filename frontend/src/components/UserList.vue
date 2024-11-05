<template>
    <div>
        <h2>User's List</h2>
        <button @click="openAddForm" id="newUser">Add New User</button>

        <!-- Search bar for filtering users -->
        <input type="text" v-model="searchQuery" placeholder="Search by name or email" />

        <!-- Show UserForm when adding or editing a user -->
        <UserForm
            v-if="showForm"
            :initialData="selectedUser"
            :isEditMode="isEditMode"
            @submit="handleFormSubmit"
            @cancel="closeForm"
        />

        <table> 
            <thead>
                <tr>
                    <th>ID</th>
                    <th>Name</th>
                    <th>Email</th>
                    <th>Age</th>
                    <th>Actions</th>
                </tr>
            </thead>
            <tbody>
                <!-- Display users from the server with pagination -->
                <tr v-for="user in users" :key="user.id"> 
                    <td>{{ user.id }}</td>
                    <td>{{ user.name }}</td>
                    <td>{{ user.email }}</td>
                    <td>{{ user.age }}</td>
                    <td>
                        <button @click="openEditForm(user)">Edit</button>
                        <button @click="deleteUser(user.id)">Delete</button>
                    </td>
                </tr>
            </tbody>
        </table>

        <!-- Pagination Controls -->
        <div id="pagination-container">
            <button @click="previousPage" :disabled="currentPage === 1">Previous</button>
            <span>Page {{ currentPage }} of {{ totalPages }}</span>
            <button @click="nextPage" :disabled="currentPage === totalPages">Next</button>
        </div>

        <!-- Messages -->
        <div v-if="message" :class="`message ${messageType}`">{{ message }}</div>
        <div v-if="successMessage" class="success-message">{{ successMessage }}</div>
    </div>
</template>

<script>
import axios from 'axios';
import UserForm from './UserForm.vue';

export default {
    name: "UserList",
    components: { UserForm },
    data() {
        return {
            users: [],
            showForm: false,
            selectedUser: null,
            isEditMode: false,
            currentPage: 1,
            pageSize: 10,
            totalUsers: 0,
            message: '',
            messageType: '',
            successMessage: '',
            searchQuery: ''
        };
    },
    methods: {
        // Fetches users from the backend
        async fetchUsers() {
            try {
                const response = await axios.get(`http://localhost:9000/users`, {
                    params: {
                        page: this.currentPage,
                        limit: this.pageSize,
                        search: this.searchQuery
                    }
                });

                // Update users and total count from server response
                this.users = response.data.users;
                this.totalUsers = response.data.totalCount;
            } catch (error) {
                this.setMessage('Error fetching users', 'error');
            }
        },
        // Navigates to the next page
        nextPage() {
            if (this.currentPage < this.totalPages) {
                this.currentPage += 1;
                this.fetchUsers();
            }
        },
        // Navigates to the previous page
        previousPage() {
            if (this.currentPage > 1) {
                this.currentPage -= 1;
                this.fetchUsers();
            }
        },
        // Opens the form for adding a new user
        openAddForm() {
            this.selectedUser = null;
            this.isEditMode = false;
            this.showForm = true;
        },
        // Opens the form with data for editing the selected user
        openEditForm(user) {
            this.selectedUser = user;
            this.isEditMode = true;
            this.showForm = true;
        },
        // Closes the add/edit form
        closeForm() {
            this.showForm = false;
        },
        // Submits the user form data to either create or update a user
        async handleFormSubmit(user) {
            try {
                if (this.isEditMode) {
                    await axios.put(`http://localhost:9000/users/${this.selectedUser.id}`, user);
                    this.setSuccessMessage('User updated successfully');
                } else {
                    await axios.post(`http://localhost:9000/users`, user);
                    this.setSuccessMessage('User added successfully');
                }
                await this.fetchUsers();
                this.closeForm();
            } catch (error) {
                this.setMessage('Error saving the user', 'error');
            }
        },
        // Deletes the selected user and refreshes the user list
        async deleteUser(id) {
            try {
                await axios.delete(`http://localhost:9000/users/${id}`);
                this.fetchUsers();
                this.setSuccessMessage('User deleted successfully');
            } catch (error) {
                this.setMessage('Error deleting the user', 'error');
            }
        },
        // Sets an error or informational message with a type
        setMessage(message, type) {
            this.message = message;
            this.messageType = type;
            setTimeout(() => {
                this.message = '';
                this.messageType = '';
            }, 4000);
        },
        // Sets a success message
        setSuccessMessage(message) {
            this.successMessage = message;
            setTimeout(() => {
                this.successMessage = '';
            }, 4000);
        }
    },
    computed: {
        // Calculates the total number of pages
        totalPages() {
            return Math.ceil(this.totalUsers / this.pageSize);
        }
    },
    watch: {
        // Resets to the first page and fetches users whenever the search query changes
        searchQuery() {
            this.currentPage = 1;
            this.fetchUsers();
        },
        // Resets to the first page and fetches users whenever the page size changes
        pageSize() {
            this.currentPage = 1;
            this.fetchUsers();
        }
    },
    // Fetches the first page of users when the component mounts
    mounted() {
        this.fetchUsers();
    }
};
</script>

<style scoped src="../assets/styles/listStyle.css"></style>
