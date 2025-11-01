import { useState } from 'react';
import { CreatePost } from '../app/groups/[id]/page.js';
// import "../styles/groupstyle.css"
import { useParams } from "next/navigation";


export function CreatePostForm({ onSubmit, onCancel }) {
  const [PostTitle, setPostTitle] = useState('');
  const [PostDescription, setPostDescription] = useState('');
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedImage, setSelectedImage] = useState(null);
  const [imagePreview, setImagePreview] = useState(null);

  const handleImageChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      setSelectedImage(file);
      // Create preview URL
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
    <div >
      <div className='group-modal-overlay' onClick={onCancel} style={{zIndex:0}} ></div>
      <div className="group-modal-content" style={{zIndex:100}}>
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
                required
              />
            </div>

            {/* post Description */}
            <div className="form-group">
              <label htmlFor="groupDescription" className="form-label">Description</label>
              <textarea
                id="postDescription"
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
                disabled={!PostTitle.trim()}
              >
                Create post
              </button>
            </div>
          </div>

        </form>
      </div>
    </div>
  );
}

export function PostCreationTrigger({ setPost }) {
  const [isModalOpen, setIsModalOpen] = useState(false);
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
      setIsModalOpen(false);
      // setShowForm(false)
      setPost(prev => {
        // console.log("posts before are :", prev);
        const exists = prev.some(p => p.id === newpost.id);
        const temp = [newpost, ...prev]
        // console.log("posts after are :", temp);
        return exists ? prev : [newpost, ...prev];
      })
    } catch (err) {
      console.error("Error creating post:", err);
    }
  };


  const userList = []
  return (
    <>
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
        />
      )}
    </>
  )
}
