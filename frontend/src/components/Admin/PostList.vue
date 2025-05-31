<template>
  <div class="post-list">
    <div class="header">
      <h1>Posts</h1>
      <router-link to="/admin/posts/new" class="new-post-button">
        New Post
      </router-link>
    </div>

    <div class="posts-grid">
      <div v-for="post in posts" :key="post.id" class="post-card">
        <div class="post-image">
          <img :src="post.cover_image" :alt="post.title" />
        </div>
        <div class="post-content">
          <h2>{{ post.title }}</h2>
          <p class="summary">{{ post.summary }}</p>
          <div class="post-meta">
            <span class="date">{{ formatDate(post.created_at) }}</span>
            <span class="read-time">{{ post.read_time }} min read</span>
          </div>
          <div class="tags">
            <span v-for="tag in post.tags" :key="tag" class="tag">
              {{ tag }}
            </span>
          </div>
          <div class="actions">
            <router-link :to="`/admin/posts/${post.id}`" class="edit-button">
              Edit
            </router-link>
            <button @click="deletePost(post.id)" class="delete-button">
              Delete
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const API_URL = 'http://localhost:8080/api'
const posts = ref([])

onMounted(async () => {
  await fetchPosts()
})

const fetchPosts = async () => {
  try {
    const response = await axios.get(`${API_URL}/posts`)
    posts.value = response.data
  } catch (error) {
    console.error('Error fetching posts:', error)
  }
}

const deletePost = async (id) => {
  if (!id) {
    console.error('ID do post não encontrado')
    return
  }

  if (confirm('Are you sure you want to delete this post?')) {
    try {
      await axios.delete(`${API_URL}/posts/${id}`)
      await fetchPosts()
    } catch (error) {
      console.error('Error deleting post:', error.response?.data || error.message)
      alert(`Error deleting post: ${error.response?.data || error.message}`)
    }
  }
}

const formatDate = (date) => {
  return new Date(date).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}
</script>

<style scoped>
.post-list {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.new-post-button {
  padding: 10px 20px;
  background: #007bff;
  color: white;
  text-decoration: none;
  border-radius: 4px;
}

.new-post-button:hover {
  background: #0056b3;
}

.posts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.post-card {
  border: 1px solid #ddd;
  border-radius: 8px;
  overflow: hidden;
  transition: transform 0.2s;
}

.post-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
}

.post-image {
  height: 200px;
  overflow: hidden;
}

.post-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.post-content {
  padding: 20px;
}

.post-content h2 {
  margin: 0 0 10px 0;
  font-size: 1.5em;
}

.summary {
  color: #666;
  margin-bottom: 15px;
}

.post-meta {
  display: flex;
  justify-content: space-between;
  color: #888;
  font-size: 0.9em;
  margin-bottom: 15px;
}

.tags {
  margin-bottom: 15px;
}

.tag {
  display: inline-block;
  padding: 3px 8px;
  margin: 0 5px 5px 0;
  background: #e9ecef;
  border-radius: 15px;
  font-size: 0.8em;
}

.actions {
  display: flex;
  gap: 10px;
}

.edit-button,
.delete-button {
  padding: 5px 15px;
  border-radius: 4px;
  cursor: pointer;
  text-decoration: none;
  text-align: center;
}

.edit-button {
  background: #28a745;
  color: white;
}

.delete-button {
  background: #dc3545;
  color: white;
  border: none;
}

.edit-button:hover {
  background: #218838;
}

.delete-button:hover {
  background: #c82333;
}
</style> 