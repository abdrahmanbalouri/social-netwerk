import "../styles/cratepost.css";
const CreatePost = ({ onCreatePost }) => {
    return (
        <div class="create-post-wrapper" onClick={onCreatePost}>
            <div class="create-post-container">
                <div class="create-post-header">
                    <div class="post-icon">
                        <i class="fa-solid fa-pen"></i>
                    </div>
                    <div class="post-input-wrapper">
                        <input type="text" class="post-input" placeholder="What's on your mind, Nicolas?" readonly />
                    </div>
                </div>
                <div class="post-divider">
                </div>
                <div class="post-actions">
                    <button class="action-button photo">
                        <i class="fa-solid fa-image"></i>
                        <span>Photo/Video</span> </button>
                    <button class="action-button video">
                        <i class="fa-solid fa-video"></i>
                        <span>Live Video</span> </button>
                    <button class="action-button feeling">
                        <i class="fa-solid fa-face-smile"></i>
                        <span>Feeling</span>
                    </button>
                </div>
            </div>
        </div>
    )
}

export default CreatePost
