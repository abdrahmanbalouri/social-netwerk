// Post.js
import Link from 'next/link';
import "../styles/post.css"


export default function Post({ post, onGetComments, ondolike }) {
console.log("dsdsdsdsd",post);



  return (
    <div className="post">
      <div className="container">
        <div className="user">
          <div className="userInfo">
            <img src={post?.profile ? `/uploads/${post.profile}` : '/assets/default.png'} alt="user" />
            <div className="details">
              <Link href={`/profile/${post.user_id}`} style={{ textDecoration: "none", color: "inherit" }}>
                <span className="name">{post.author}</span>
              </Link>
              <span className="date">{new Date(post.created_at).toLocaleString()}</span>
            </div>
          </div>
        </div>
        <div className="content1">
          <h3>
            <p style={{ color: "#5271ff" }}>{post.title}</p>
            <strong className="content2">{post.content}</strong>
          </h3>

          {/* Check if image_path exists */}
          {post.image_path && (
            post.image_path.match(/\.(jpeg|jpg|png|gif|webp)$/i) ? (
              <img src={`/${post.image_path}`} alt="Post content" className="post-media-image" />
            ) : post.image_path.match(/\.(mp4|webm|ogg|mov)$/i) ? (
              <div className="post-media-wrapper">
                <video controls className="post-media-video">
                  <source src={`/${post.image_path}`} type={`video/${post.image_path.split('.').pop()}`} />
                  Your browser does not support the video tag.
                </video>
              </div>
            ) : null
          )}


        </div>
        <div className="infoo">
          <div className="item" onClick={() => ondolike(post.id)} >
            <i
              className={post.liked_by_user ? "fa-solid fa-heart" : "fa-regular fa-heart"}
              style={post.liked_by_user ? { color: "red" } : {}}
            />
            {post.like || 0} Likes
          </div>
          <div className="item" onClick={() => {
            onGetComments(post)


          }

          }>
            <i className="fa-solid fa-comment"></i> {post.comments_count || 0} Comments
          </div>
        </div>
      </div>
    </div>
  );
}