// Post.js
import Link from 'next/link';

export default function Post({ post, onGetComments, ondolike }) {
  


  return (
    <div className="post">
      <div className="container">
        <div className="user">
          <div className="userInfo">
            <img src={`/uploads/${post.profile}` || '/avatar.png'} alt="user" />
            <div className="details">
              <Link href={`/profile/${post.user_id}`} style={{ textDecoration: "none", color: "inherit" }}>
                <span className="name">{post.author}</span>
              </Link>
              <span className="date">{new Date(post.created_at).toLocaleString()}</span>
            </div>
          </div>
        </div>
        <div className="content">
          <h3>
            <p style={{ color: "#5271ff" }}>{post.title}</p>
            {post.content}
          </h3>
          {post.image_path && <img src={`/${post.image_path}`} alt="Post content" />}
        </div>
        <div className="info">
          <div className="item" onClick={() => ondolike(post.id)} >
            <i
              className={post.liked_by_user ? "fa-solid fa-heart" : "fa-regular fa-heart"}
              style={post.liked_by_user ? { color: "red" } : {}}
            />
            {post.like || 0} Likes
          </div>
          <div className="item" onClick={() => onGetComments(post)}>
            <i className="fa-solid fa-comment"></i> {post.comments_count || 0} Comments
          </div>
        </div>
      </div>
    </div>
  );
}