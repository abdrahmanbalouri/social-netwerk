import { useState } from 'react';
import { CreatePost } from '../app/groups/[id]/page.js';
import "../styles/groupstyle.css"
import { useParams } from "next/navigation";


export function CreatePostForm({ onSubmit, onCancel }) {
    const [PostTitle, setPostTitle] = useState('');
    const [PostDescription, setPostDescription] = useState('');
    const [searchQuery, setSearchQuery] = useState('');

    const handleSubmit = (e) => {
        e.preventDefault();

        const postData = {
            title: PostTitle,
            description: PostDescription,
        };

        onSubmit(postData);

        // Reset form after submission
        setPostTitle('');
        setPostDescription('');
        setSearchQuery('');
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


    const userList = []
    return (
        <>
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

        </>
    )
}

