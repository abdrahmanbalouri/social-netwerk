import "../styles/cratepost.css";
const CreatePost = ({ onCreatePost }) => {
    return (
        <div className="create-post-wrapper" onClick={onCreatePost}>
            <div className="create-post-container">
                <div className="create-post-header">
                    <div className="post-icon">
                        <i className="fa-solid fa-pen"></i>
                    </div>
                    <div className="post-input-wrapper">
                        <input type="text" className="post-input" placeholder="What's on your mind, Nicolas?" readOnly />
                    </div>
                </div>
                <div className="post-divider">
                </div>
                <div className="post-actions">
                    <button className="action-button photo">
                        <i className="fa-solid fa-image"></i>
                        <span>Photo/Video</span> </button>
                    <button className="action-button video">
                        <i className="fa-solid fa-video"></i>
                        <span>Live Video</span> </button>
                    <button className="action-button feeling">
                        <i className="fa-solid fa-face-smile"></i>
                        <span>Feeling</span>
                    </button>
                </div>
            </div>
        </div>
    )
}

export default CreatePost
