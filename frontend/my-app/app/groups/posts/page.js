"use client";
import Navbar from "../../../components/Navbar.js"
import Post from "../../../components/Post.js";
import { useEffect, useState } from "react";
import "./page.css"


export default function () {
    return (
        <>
            <Navbar />
            <AllPosts />
        </>
    )
}

function AllPosts() {

    const [posts, setPost] = useState(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const grpID = localStorage.getItem('selectedGroup')
        if (!grpID) return
        fetch('http://localhost:8080/group/fetchPosts', {
            method: 'POST',
            credentials: 'include',
            body: JSON.stringify({
                grpId: grpID
            }),
            credentials: 'include',
        })
            .then(res => res.json())
            .then(data => {
                console.log("data is :", data);
                setPost(data);
                setLoading(false);
            })
            .catch(error => {
                console.error("Failed to fetch posts for data:", error);
                setLoading(false);
            });
    }, [])
    if (!posts) {
        return <div>There is no post yetttttt.</div>;
    }
    console.log("posts are :", posts);
    return (
        <>
            <CreatePost />
            <div className="posts-list">
                {posts.map((post) => (
                        <Post
                            key={post.id}
                            post={post}
                            onGetComments={GetComments}
                            ondolike={AddLike}
                        />
                        ))}
            </div>
        </>
    )
}
function CreatePost() {
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [postContent, setPostContent] = useState('');
    const [userName, setUserName] = useState('Ussra'); //example

    const handlePostClick = () => {
        setIsModalOpen(true);
    };

    const handleCloseModal = () => {
        setIsModalOpen(false);
        setPostContent('');
    };

    const handlePostSubmit = () => {
        if (postContent.trim()) {
            // Send post to backend
            fetch('http://localhost:8080/group/addPost', {
                method: 'POST',
                credentials: 'include',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    content: postContent,
                    groupId: localStorage.getItem('selectedGroupId'),
                }),
            })
                .then(res => res.json())
                .then(data => {
                    console.log("Post created:", data);
                    setPostContent('');
                    setIsModalOpen(false);
                    // Refresh posts or update state here
                })
                .catch(error => console.error("Failed to create post:", error));
        }
    };

    return (
        <>
            <div className="create-post-container">
                <div className="create-post-card">
                    <div className="create-post-header">
                        <div className="user-avatar">
                            <span>{userName}</span>
                        </div>
                        <input
                            type="text"
                            placeholder="What's on your mind?"
                            className="post-input"
                            onClick={handlePostClick}
                            readOnly
                        />
                    </div>
                </div>
            </div>

            {isModalOpen && (
                <div className="modal-overlay" onClick={handleCloseModal}>
                    <div className="modal-content" onClick={(e) => e.stopPropagation()}>
                        <div className="modal-header">
                            <h2>Create Post</h2>
                            <button className="close-button" onClick={handleCloseModal}>✕</button>
                        </div>
                        <div className="modal-body">
                            <div className="modal-user-info">
                                <div className="user-avatar">
                                    <span>{userName}</span>
                                </div>
                                <div>
                                    <p className="user-name">{userName}</p>
                                    <p className="user-status">Public</p>
                                </div>
                            </div>
                            <textarea
                                value={postContent}
                                onChange={(e) => setPostContent(e.target.value)}
                                placeholder="What's on your mind?"
                                className="modal-textarea"
                            />
                        </div>
                        <div className="modal-footer">
                            <button className="cancel-button" onClick={handleCloseModal}>Cancel</button>
                            <button
                                className="submit-button"
                                onClick={handlePostSubmit}
                                disabled={!postContent.trim()}
                            >
                                Post
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </>
    );
}

function GetComments(post) {
    // setSelectedPost({
    //     id: post.id,
    //     title: post.title || post.post_title || "Post"
    // });
}

function AddLike() {

}
