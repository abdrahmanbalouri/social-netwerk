import { useState } from "react";
import { CreatePost } from "../app/groups/[id]/page.js";
import "../styles/groupstyle.css";
import { useParams } from "next/navigation";

export function CreatePostForm({ onSubmit, onCancel }) {
  const [PostTitle, setPostTitle] = useState("");
  const [PostDescription, setPostDescription] = useState("");
  const [searchQuery, setSearchQuery] = useState("");

  const handleSubmit = (e) => {
    e.preventDefault();

    const postData = {
      title: PostTitle,
      description: PostDescription,
    };

    onSubmit(postData);

    // Reset form after submission
    setPostTitle("");
    setPostDescription("");
    setSearchQuery("");
  };

  return (
    <div className="group-modal-overlay">
      <div className="group-modal-content">
        <div className="group-modal-header">
          <h1 className="group-modal-title">Create a New Post</h1>
        </div>

        <form onSubmit={handleSubmit}>
          {/* post Title */}
          <div className="group-modal-form">
            <div className="form-group">
              <label htmlFor="postTitle" className="form-label">
                Post Title
              </label>
              <input
                id="postTitle"
                type="text"
                value={PostTitle}
                onChange={(e) => setPostTitle(e.target.value)}
                placeholder="Enter post title"
                className="form-input"
                required
                autoFocus
              />
            </div>

            {/* post Description */}
            <div className="form-group">
              <label htmlFor="postContent" className="form-label">
                Content
              </label>
              <textarea
                id="postDescription"
                value={PostDescription}
                onChange={(e) => setPostDescription(e.target.value)}
                placeholder="What's on your mind?"
                className="form-textarea"
                required
              />
            </div>
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
        </form>
      </div>
    </div>
  );
}

export function PostCreationTrigger() {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const { id } = useParams();

  const handlePostClick = () => {
    // setShowNoGroups(false);
    setIsModalOpen(true);
  };
  const handleCancel = () => {
    setIsModalOpen(false);
    // setShowNoGroups(true);
  };
  const handleSubmit = async (formData) => {
    try {
      await CreatePost(id, formData);
      setIsModalOpen(false);
    } catch (err) {
      console.error("Error creating post:", err);
    }
  };

  const userList = [];
  return (
    <>
      <div className="create-card">
        <div className="create-card-inner">
          <div className="avatar">U</div>
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
  );
}
