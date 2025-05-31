<template>
  <div class="post-editor">
    <div class="editor-header">
      <input
        v-model="post.title"
        type="text"
        placeholder="Post Title"
        class="title-input"
      />
      <input
        v-model="post.summary"
        type="text"
        placeholder="Post Summary"
        class="summary-input"
      />
      <div class="cover-image">
        <div class="cover-image-preview" v-if="post.cover_image">
          <img :src="post.cover_image" alt="Cover" />
          <button @click="post.cover_image = ''" class="remove-image">×</button>
        </div>
        <div class="cover-image-input">
          <input
            type="text"
            v-model="post.cover_image"
            placeholder="Cover Image URL"
          />
          <button @click="insertUnsplashImage" class="unsplash-button">
            Add Unsplash Image
          </button>
        </div>
      </div>
    </div>

    <div class="editor-toolbar">
      <button @click="insertImage">Image</button>
      <button @click="insertUnsplashImage">Unsplash Image</button>
      <button @click="insertVideo">Video</button>
      <button @click="insertEmbed">Embed</button>
      <button @click="insertCodeBlock">Code Block</button>
      <button @click="insertNewPart">New Part</button>
    </div>

    <div class="editor-content">
      <textarea
        v-model="post.content"
        placeholder="Write your post content here..."
        @keyup="handleContentChange"
      ></textarea>
    </div>

    <div class="tags-section">
      <input
        v-model="newTag"
        @keyup.enter="addTag"
        placeholder="Add tags (press Enter)"
      />
      <div class="tags-list">
        <span v-for="(tag, index) in post.tags" :key="index" class="tag">
          {{ tag }}
          <button @click="removeTag(index)">×</button>
        </span>
      </div>
    </div>

    <div class="editor-footer">
      <button @click="savePost" class="save-button">Save Post</button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'

const route = useRoute()
const API_URL = 'http://localhost:8080/api'

const post = ref({
  title: '',
  summary: '',
  content: '',
  cover_image: '',
  tags: [],
  author: {
    name: 'Admin',
    image: '',
    role: 'Author'
  }
})

const newTag = ref('')

onMounted(async () => {
  if (route.params.id) {
    try {
      const response = await axios.get(`${API_URL}/posts/${route.params.id}`)
      post.value = response.data
    } catch (error) {
      console.error('Error fetching post:', error)
    }
  }
})

const addTag = () => {
  if (newTag.value.trim()) {
    post.value.tags.push(newTag.value.trim())
    newTag.value = ''
  }
}

const removeTag = (index) => {
  post.value.tags.splice(index, 1)
}

const insertImage = () => {
  const url = prompt('Enter image URL:')
  if (url) {
    insertAtCursor(`<img src="${url}" alt="Image" />`)
  }
}

const insertUnsplashImage = () => {
  const query = prompt('Enter search query for Unsplash:')
  if (query) {
    const url = `https://source.unsplash.com/random/800x600/?${query}`
    post.value.cover_image = url
  }
}

const insertVideo = () => {
  const url = prompt('Enter video URL:')
  if (url) {
    insertAtCursor(`<video src="${url}" controls></video>`)
  }
}

const insertEmbed = () => {
  const url = prompt('Enter embed URL:')
  if (url) {
    insertAtCursor(`<iframe src="${url}" frameborder="0" allowfullscreen></iframe>`)
  }
}

const insertCodeBlock = () => {
  const code = prompt('Enter code:')
  if (code) {
    insertAtCursor(`<pre><code>${code}</code></pre>`)
  }
}

const insertNewPart = () => {
  insertAtCursor('\n\n---\n\n')
}

const insertAtCursor = (text) => {
  const textarea = document.querySelector('.editor-content textarea')
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  post.value.content = post.value.content.substring(0, start) + text + post.value.content.substring(end)
}

const handleContentChange = () => {
  // Calculate reading time (assuming 200 words per minute)
  const words = post.value.content.trim().split(/\s+/).length
  post.value.readTime = Math.ceil(words / 200)
}

const savePost = async () => {
  try {
    // Validar campos obrigatórios
    if (!post.value.title) {
      alert('Please enter a title')
      return
    }

    if (!post.value.content) {
      alert('Please enter some content')
      return
    }

    console.log('Enviando post:', post.value)
    if (route.params.id) {
      await axios.put(`${API_URL}/posts/${route.params.id}`, post.value)
    } else {
      const response = await axios.post(`${API_URL}/posts`, post.value)
      console.log('Resposta do servidor:', response.data)
    }
    alert('Post saved successfully!')
  } catch (error) {
    console.error('Error saving post:', error)
    if (error.response) {
      console.error('Status:', error.response.status)
      console.error('Data:', error.response.data)
      console.error('Headers:', error.response.headers)
      alert(`Error saving post: ${JSON.stringify(error.response.data)}`)
    } else if (error.request) {
      console.error('Request:', error.request)
      alert('Error saving post: No response from server')
    } else {
      console.error('Error:', error.message)
      alert(`Error saving post: ${error.message}`)
    }
  }
}
</script>

<style scoped>
.post-editor {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.editor-header {
  margin-bottom: 20px;
}

.title-input {
  font-size: 2em;
  width: 100%;
  margin-bottom: 10px;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.summary-input {
  width: 100%;
  margin-bottom: 10px;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.cover-image {
  margin-bottom: 20px;
}

.cover-image-preview {
  position: relative;
  margin-bottom: 10px;
}

.cover-image-preview img {
  max-width: 100%;
  max-height: 300px;
  border-radius: 4px;
}

.remove-image {
  position: absolute;
  top: 10px;
  right: 10px;
  background: rgba(255, 255, 255, 0.8);
  border: none;
  border-radius: 50%;
  width: 30px;
  height: 30px;
  font-size: 20px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.cover-image-input {
  display: flex;
  gap: 10px;
}

.cover-image-input input {
  flex: 1;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.unsplash-button {
  padding: 10px 20px;
  background: #28a745;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.unsplash-button:hover {
  background: #218838;
}

.editor-toolbar {
  margin-bottom: 20px;
  padding: 10px;
  background: #f5f5f5;
  border-radius: 4px;
}

.editor-toolbar button {
  margin-right: 10px;
  padding: 5px 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: white;
  cursor: pointer;
}

.editor-content textarea {
  width: 100%;
  min-height: 400px;
  padding: 20px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 16px;
  line-height: 1.6;
}

.tags-section {
  margin: 20px 0;
}

.tags-list {
  margin-top: 10px;
}

.tag {
  display: inline-block;
  padding: 5px 10px;
  margin: 0 5px 5px 0;
  background: #e9ecef;
  border-radius: 20px;
}

.tag button {
  margin-left: 5px;
  border: none;
  background: none;
  cursor: pointer;
}

.editor-footer {
  margin-top: 20px;
  text-align: right;
}

.save-button {
  padding: 10px 20px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.save-button:hover {
  background: #0056b3;
}
</style> 