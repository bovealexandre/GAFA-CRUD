<template>
    <i class="fas fa-trash-alt cursor-pointer" style="color: hsl(0,50%,50%)" @click="deleteItem(props.id)"></i>
    <div style="display:none;">{{props.id}}</div>
</template>

<script lang="ts">
    export default {
        props: {
            id: {
                type: String,
                required: true,
            }
        },
        setup(props) {
            props
            const deleteItem = (id: string) => {
                return fetch('http://localhost:5000/api/todos/' + id, {
                    method: 'DELETE',
                }).then(response => response.json())
            }
            return {
                props,
                deleteItem,
            };
        },
    };
</script>