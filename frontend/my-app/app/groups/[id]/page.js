"use client";
import Navbar from "../../../components/Navbar.js";
import { GroupPostChat } from "../../../components/groupPostCat.js";
import Post from "../../../components/Post.js";
import { useEffect, useState } from "react";
import "../../../styles/groupstyle.css";
import { useParams, useRouter } from "next/navigation";
import { PostCreationTrigger } from "../../../components/cretaePostGroup.js";
import LeftBar from "../../../components/LeftBar.js";
import RightBarGroup from "../../../components/RightBarGroups.js";
import { useDarkMode } from "../../../context/darkMod.js";
import Comment from "../../../components/coment.js";

// Global sendRequest (can be moved to a service file later)
async function sendRequest(invitedUserID, grpID) {
  console.log("invited user id is : ", invitedUserID);
  console.log("group id is : ", grpID);

  try {
    const res = await fetch(`http://localhost:8080/group/invitation/${grpID}`, {
      method: 'POST',
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        invitedUsers: [invitedUserID],
      }),
    });

    const temp = await res.json();
    console.log("temp is :", temp);
    return temp;
  } catch (error) {
    console.error("error sending invitation to the user :", error);
    return { error: error.message };
  }
}

// Close comments modal
function closeComments(setShowComments, setSelectedPost, setComment) {
  setShowComments(false);
  setSelectedPost(null);
  setComment([]);
}

// Main Page Component
export default function GroupPage() {
  const { darkMode } = useDarkMode();
  const router = useRouter();

  return (
    <div className={darkMode ? "theme-dark" : "theme-light"}>
      <Navbar />
      <main className="content">
        <LeftBar showSidebar={true} />
        <GroupPostChat />
        <RightBarGroup onClick={sendRequest} />
      </main>
    </div>
  );
}

// AllPosts Component
export function AllPosts() {
  const [posts, setPost] = useState([]);
  const [loading, setLoading] = useState(true);
  const [loadingComment, setLoadingComment] = useState(false);
  const [selectedPost, setSelectedPost] = useState(null);
  const [showComments, setShowComments] = useState(false);
  const [comment, setComment] = useState([]);
  
  const params = useParams();
  const router = useRouter();
  const grpID = params.id;

  useEffect(() => {
    if (!grpID) {
      setLoading(false);
      return;
    }

    fetch(`http://localhost:8080/group/fetchPosts/${grpID}`, {
      method: 'GET',
      credentials: 'include',
    })
      .then(res => {
        if (!res.ok) throw new Error('Failed to fetch posts');
        return res.json();
      })
      .then(data => {
        setPost(data || []);
        setLoading(false);
      })
      .catch(error => {
        console.error("Failed to fetch posts:", error);
        setLoading(false);
      });
  }, [grpID]);

  // Handle Like
  const handleLike = async (postId) => {
    
    try {
      const res = await fetch(`http://localhost:8080/api/like/${postId}`, {
        method: "POST",
        credentials: "include",
      });
      const response = await res.json();
      console.log(response);
      

      if (response.error) {
        if (response.error === "Unauthorized") {
          router.push("/login");
        }
        return;
      }

      const updatedPost = await fetchPostById(postId);
      //console.log(updatedPost);
      
      
      if (updatedPost) {
        setPost(prevPosts =>
          prevPosts.map(p => p.id === updatedPost.id ? updatedPost : p)
        );
      }
    } catch (err) {
      console.error("Error liking post:", err);
    }
  };

  const fetchPostById = async (postID) => {
    try {
      const res = await fetch(`http://localhost:8080/group/updatepost/${postID}`, {
        method: "GET",
        credentials: "include",
      });
      const data = await res.json();
      console.log(data);
      
      return data.error ? null : data;
    } catch (err) {
      console.error("Error fetching post:", err);
      return null;
    }
  };

  const getComments = async (post) => {
    setLoadingComment(true);
    try {
      setSelectedPost({
        id: post.id,
        title: post.title || post.post_title || "Post",
        image_path: post.image_path,
        content: post.content,
        author: post.first_name + " " + post.last_name,
      });
      setShowComments(true);

      const res = await fetch(`http://localhost:8080/group/fetchComments/${post.id}`, {
        method: "GET",
        credentials: "include",
      });

      if (!res.ok) {
        const errorText = await res.text();
        console.error("Failed to fetch comments:", res.status, errorText);
        return false;
      }

      const data = await res.json();
      setComment(data || []);
      return data?.[0]?.id || null;
    } catch (err) {
      console.error("Error fetching comments:", err);
      return false;
    } finally {
      setLoadingComment(false);
    }
  };

  return (
    <div>
      <PostCreationTrigger setPost={setPost} groupId={grpID} />
      
      {loading ? (
        <div>Loading posts...</div>
      ) : posts.length === 0 ? (
        <div>There is no post yet.</div>
      ) : (
        <div className="posts-list">
          {posts.map((post) => (
            <Post
              key={post.id}
              post={post}
              onGetComments={() => getComments(post)}
              ondolike={handleLike}
            />
          ))}

          <Comment
            comments={comment}
            isOpen={showComments}
            onClose={() => closeComments(setShowComments, setSelectedPost, setComment)}
            postId={selectedPost?.id}
            postTitle={selectedPost?.title}
            ongetcomment={getComments}
            post={selectedPost}
            loading={loadingComment}
          />
        </div>
      )}
    </div>
  );
}

// Create Post (External function)
export async function CreatePost(groupId, formData) {
  try {
    const res = await fetch(`http://localhost:8080/group/addPost/${groupId}`, {
      method: "POST",
      credentials: "include",
      body: formData,
    });

    const text = await res.text();
    console.log("Raw response:", text);

    if (!res.ok) {
      throw new Error(`Failed to create post: ${res.status} ${text}`);
    }

    return JSON.parse(text);
  } catch (error) {
    console.error("CreatePost error:", error);
    throw error;
  }
}