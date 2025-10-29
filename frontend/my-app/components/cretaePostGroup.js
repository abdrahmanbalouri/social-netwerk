import { useState } from 'react';
import { CreatePost } from '../app/groups/[id]/page.js';
import "../styles/groupstyle.css"
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

        const postData = {
            title: PostTitle,
            description: PostDescription,
            image: selectedImage,
        };
        onSubmit(postData);
        // Reset form after submission
        setPostTitle('');
        setPostDescription('');
        setSearchQuery('');
        setSelectedImage(null);
        setImagePreview(null);
    };

    return (
        <div className="create-post-form-container">
            <div className="create-post-form-card">
                <h1 className="form-title">Create a New post</h1>

                <form onSubmit={handleSubmit}>
                    {/* post Title */}
                    <div className="form-field">
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
                    <div className="form-field">
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
                            accept="image/*"
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
                    <div className="form-actions">
                        <button
                            type="button"
                            className="cancel-post-btn"
                            onClick={onCancel}
                        >
                            Cancel
                        </button>
                        <button
                            type="submit"
                            className="submit-post-btn"
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
                console.log("posts before are :", prev);
                const exists = prev.some(p => p.id === newpost.id);
                const temp = [newpost, ...prev]
                console.log("posts after are :", temp);
                return exists ? prev : [newpost, ...prev];
            })
        } catch (err) {
            console.error("Error creating post:", err);
        }
    };


    const userList = []
    return (
        <div>
            <div className="create-post-container">
                <div className="create-post-card">
                    <div className="create-post-header">
                        <div className="user-avatar">
                            <span>User</span>
                        </div>
                        <input
                            type="text"
                            placeholder="Do you want to create a post? Just click here!"
                            className="post-input"
                            onClick={handlePostClick}
                            readOnly
                        />
                    </div>
                </div>
            </div>
            {isModalOpen && (
                <CreatePostForm
                    users={userList}
                    onSubmit={handleSubmit}
                    onCancel={handleCancel}
                />
            )}

        </div>
    )
}

