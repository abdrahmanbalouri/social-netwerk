"use client";
import { useState } from 'react';
import { CreatePost } from '../app/groups/[id]/page.js';
import { Toaster, toast } from "sonner"
import { useParams, useRouter } from "next/navigation";

export function CreatePostForm({ onSubmit, onCancel, err }) {
  const [PostTitle, setPostTitle] = useState('');
  const [PostDescription, setPostDescription] = useState('');
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedImage, setSelectedImage] = useState(null);
  const [imagePreview, setImagePreview] = useState(null);

  const handleImageChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      setSelectedImage(file);
      const reader = new FileReader();
      reader.onloadend = () => {
        setImagePreview(reader.result);
      };
      reader.readAsDataURL(file);
    }
  };

  const removeImage = () => {
    setSelectedImage(null);
    setImagePreview(null);
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    const formData = new FormData();
    formData.append("title", PostTitle);
    formData.append("description", PostDescription);
    formData.append("image", selectedImage);

    onSubmit(formData);
    setPostTitle("");
    setPostDescription("");
    setSearchQuery("");
    setSelectedImage(null);
    setImagePreview(null);
  };


  return (
    <div className='group-modal-overlay'>
      <div className='group-modal-overlay1' onClick={onCancel}></div>
      <div className="group-modal-content">
        <div className="group-modal-header">
          <h1 className="group-modal-title">Create a New Post</h1>
        </div>
        <form onSubmit={handleSubmit}>
          {/* post Title */}
          <div className="group-modal-form">

            <div className="form-group">
              <label htmlFor="groupTitle" className="form-label">post Title</label>
              <input
                id="postTitle"
                type="text"
                value={PostTitle}
                onChange={(e) => setPostTitle(e.target.value)}
                placeholder="Enter post title"
                className="form-input"
              />
            </div>

            {/* post Description */}
            <div className="form-group">
              <label htmlFor="groupDescription" className="form-label">Description</label>
              <textarea
                id="postDescription"
                type="text"
                value={PostDescription}
                onChange={(e) => setPostDescription(e.target.value)}
                placeholder="What's your post about?"
                className="form-textarea"
                rows="4"
              />
            </div>

            {/* Image Upload */}
            <div className="form-field">
              <label htmlFor="postImage" className="form-label">Image (optional)</label>
              <input
                id="postImage"
                type="file"
                onChange={handleImageChange}
                className="form-input"
                accept="image/*,video/*"
              />
              {imagePreview && (
                <div className="image-preview-container" style={{ marginTop: '10px' }}>
                  <img
                    src={imagePreview}
                    alt="Preview"
                    style={{ maxWidth: '200px', maxHeight: '200px', borderRadius: '8px' }}
                  />
                  <button
                    type="button"
                    onClick={removeImage}
                    className="remove-image-btn"
                    style={{ marginLeft: '10px' }}
                  >
                    Remove
                  </button>
                </div>
              )}
            </div>

            {/* Submit Button */}
            <div className="group-modal-actions">
              <button type="button" className="cancel-button" onClick={onCancel}>
                Cancel
              </button>
              <button
                type="submit"
                className="submit-button"
                disabled={!PostTitle.trim() && !PostDescription.trim() && !selectedImage}
              >
                Create post
              </button>
            </div>
            {err}
          </div>

        </form>
      </div>
    </div>
  );
}

export function PostCreationTrigger({ setPost }) {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [err, seterr] = useState("")
  const router = useRouter()
  // const [posts, setPost] = useState([])
  const { id } = useParams();

  const handlePostClick = () => {
    setIsModalOpen(true);
  };
  const handleCancel = () => {
    setIsModalOpen(false);
  };
  const handleSubmit = async (formData) => {
    try {
      const newpost = await CreatePost(id, formData);
      toast.success("post created successfully!");
      if (newpost.error) {
        if (newpost.error == "Authentication required") {
          router.push('/login')
        }
        seterr(newpost.error)
        return
      }

      setIsModalOpen(false);
      setPost(prev => {
        const exists = prev.some(p => p.id === newpost.id);
        const temp = [newpost, ...prev]
        return exists ? prev : [newpost, ...prev];
      })
    } catch (err) {
      toast.error(err.message);
    }
  };


  const userList = []
  return (
    <div>
      <Toaster position="bottom-right" richColors />
      <div className="create-card">
        <div className="create-card-inner">
          <div className="group-avatar">U</div>
          <input
            type="text"
            placeholder="What's on your mind?"
            className="create-input"
            onClick={handlePostClick}
            readOnly
          />

        </div>
      </div>
      {isModalOpen && (
        <CreatePostForm
          users={userList}
          onSubmit={handleSubmit}
          onCancel={handleCancel}
          err={err}
        />
      )}
    </div>
  )
}
