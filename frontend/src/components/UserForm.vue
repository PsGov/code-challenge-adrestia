<template>
    <div>
        <h3>{{ isEditMode ? "Edit User" : "Add New User"}}</h3>
        <form>
            <div>
                <label for="name">Name: </label>
                <input type="text" id="name" v-model="formData.name" required />
            </div>
            <div>
                <label for="email">Email: </label>
                <input type="email" id="email" v-model="formData.email" required />
                <p v-if="emailError" class="err">{{ emailError }}</p>
            </div>
            <div>
                <label for="age">Age: </label>
                <input type="number" id="age" v-model="formData.age" required min="1" />
                <p v-if="ageError" class="err">{{ ageError }}</p>
            </div>
            <button type="button" @click="handleSubmit">{{ isEditMode ? "Update User" : "Add User" }}</button>
            <button type="button" @click="$emit('cancel')">Cancel</button>
        </form>
    </div>
</template>

<script>
export default {
    name: "UserForm",
    props: {
        // initialData is used to pre-fill the form when editing an existing user
        initialData: Object,
        // isEditMode indicates if we're adding a new user or editing an existing one
        isEditMode: {
            type: Boolean,
            default: false
        }
    },
    data() {
        return{
            // formData holds the form's input values, with defaults set from initialData if provided
            formData: {
                name: this.initialData?.name || '',
                email: this.initialData?.email || '',
                age: this.initialData?.age || 1
            },
            // error messages for email and age
            emailError: '',
            ageError: ''
        };
    },
    methods: {
        handleSubmit() {
            // Reset errors
            this.emailError = '';
            this.ageError = '';

            // Basic email validation
            const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            if(!emailRegex.test(this.formData.email)) {
                this.emailError = 'Invalid email format';
                return;
            }

            // Check if age is >0
            if(this.formData.age <=0 ) {
                this.ageError = 'Age must be more that 0';
                return;
            }

            // Check if all required fields are filled
            if (!this.formData.name || !this.formData.email || !this.formData.age) {
                console.warn("Incomplete data on submit:", this.formData);
                return;
            }

            this.$emit('submit', { ...this.formData });
        }

    },
    watch: {
        // Watch for changes in initialData to update the form fields dynamically
        initialData: {
            immediate: true,
            handler(newData) {
                if (newData) {
                    this.formData = { ...newData };
                } else {
                    // Reset form to default values if no initialData is provided
                    this.formData = { name: '', email: '', age: 1 };
                }
            }
        }
    },
};

</script>
<style scoped src="../assets/styles/formStyle.css"></style>
