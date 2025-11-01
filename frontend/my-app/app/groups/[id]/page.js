"use client";
import "../../../styles/events.css"
import Navbar from "../../../components/Navbar.js";
import { GroupPostChat } from "../../../components/groupPostCat.js";
import Post from "../../../components/Post.js";
import { useEffect, useState, useRef } from "react";
import "../../../styles/groupstyle.css";
import { useParams, useRouter } from "next/navigation";
import { PostCreationTrigger } from "../../../components/cretaePostGroup.js";
import LeftBar from "../../../components/LeftBar.js";
import RightBarGroup from "../../../components/RightBarGroups.js";
import { useDarkMode } from "../../../context/darkMod.js";
import Comment from "../../../components/coment.js";
import { middleware } from "../../../middleware/middelware.js";
import EventCard from "../../../components/EventCard.js";

// Global sendRequest (can be moved to a service file later)
async function sendRequest(invitedUserID, grpID) {

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

// Main Page Component
export default function GroupPage() {
  const { darkMode } = useDarkMode();
  const router = useRouter();
  const params = useParams();
  const grpID = params.id;
  useEffect(() => {

    const checkAuth = async () => {
      const auth = await middleware();
      if (!auth) {

        router.push("/login");

      }
    }
    checkAuth();
  }, [grpID])

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




// events 


export function Events() {
  const [ShowEventForm, SetShowEventForm] = useState(true)
  const params = useParams();

  const grpID = params.id;

  const [EventList, setEventList] = useState([])

  useEffect(() => {
    fetchEvents();
  }, []);

  async function fetchEvents() {
    try {
      const response = await fetch(`http://localhost:8080/api/getEvents/${grpID}`, {
        method: 'GET',
        credentials: 'include',
      });

      if (!response.ok) {
        throw new Error('Failed to fetch events');
      }

      const data = await response.json();

      console.log('Fetched events:', data);
      setEventList(data);
    } catch (error) {
      console.error('Error fetching events:', error);
    }
  }

  function showEvent() {
    SetShowEventForm(prev => !prev)
  }



  async function goingEvent(status, eventID) {

    if ((!status || !eventID) || (status !== "going" && status !== "notGoing")) {
      console.error('Status and Event ID are required');
      return;
    }
    try {
      const response = await fetch(`http://localhost:8080/api/event/action/${grpID}`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ status, eventID }),
      });

      if (!response.ok) {
        throw new Error('Failed to RSVP to event');
      }

      await fetchEvents();
    } catch (error) {
      console.error('Error RSVPing to event:', error);
    }
  }


  return (
    <>
      {ShowEventForm ? (
        <div>
          <button onClick={showEvent} className="Showbutton">
            create event
          </button>
        </div>
      ) : (
        <div className="privacy-overlay">

          <div className="backLayer" onClick={showEvent}></div>
          <EventForm ShowEventForm={ShowEventForm} closeForm={showEvent}
            fetchEvents={fetchEvents}
          />
        </div>
      )}

      {EventList?.length === 0 ? (
        <div>No events yet. Be the first to create one!</div>
      ) : (
        EventList.map((ev) => (
          <EventCard key={ev.id} ev={ev} goingEvent={goingEvent} />
        ))
      )}





    </>
  )
}



// AllPosts Component
export function AllPosts() {
  const [posts, setPost] = useState([]);
  const [loading, setLoading] = useState(true);
  const [loadingcomment, setLoadingComment] = useState(false);
  const [selectedPost, setSelectedPost] = useState(null);
  const [showComments, setShowComments] = useState(false);
  const [comment, setComment] = useState([]);
  const params = useParams();
  const router = useRouter();
  const grpID = params.id;
  const [toast, setToast] = useState(null);
  const offsetcomment = useRef(0);
  const showToast = (message, type = "error", duration = 3000) => {
    setToast({ message, type });
    setTimeout(() => {
      setToast(null);
    }, duration);
  };
  function closeComments() {
    offsetcomment.current = 0
    setShowComments(false);
    setSelectedPost(null);
    setComment([]);
  }

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

        if (data.error) {
          showToast(data.error)
          setLoading(false);
          router.push("/login");
          return
        }

        setPost(data || []);
        setLoading(false);

      })
      .catch(error => {
        console.error("Failed to fetch posts:", error);
        setLoading(false);
      });

  }, [grpID]);
  async function refreshComments(commentID) {
    if (!selectedPost?.id) return;

    try {
      const res = await fetch(`http://localhost:8080/group/getlastcomment/${commentID}/${grpID}`, {
        method: "GET",
        credentials: "include",
      });

      const data = await res.json();
      if (data.error) {
        showToast(data.error)
        return
      }

      let newcomment = [];

      if (Array.isArray(data)) {
        newcomment = data;
      } else if (data && data.newcomment && Array.isArray(data.newcomment)) {
        newcomment = data.newcomment;
      } else if (data) {
        newcomment = [data];
      }


      setComment([...newcomment, ...comment]);
      offsetcomment.current++


      const potsreplace = await fetchPostById(selectedPost.id)
      for (let i = 0; i < posts.length; i++) {
        if (posts[i].id == selectedPost.id) {
          setPost([
            ...posts.slice(0, i),
            potsreplace,
            ...posts.slice(i + 1)
          ]);
          break
        }
      }

    } catch (err) {
      console.error("Error refreshing comments:", err);
    }
  }

  // Handle Like
  const handleLike = async (postId) => {

    try {
      const res = await fetch(`http://localhost:8080/api/like/${postId}/${grpID}`, {
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
      const res = await fetch(`http://localhost:8080/group/updatepost/${postID}/${grpID}`, {
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

  async function GetComments(post) {
    setLoadingComment(true)

    try {
      setSelectedPost({
        id: post.id,
        title: post.title || post.post_title || "Post",
        image_path: post.image_path,
        content: post.content,
        author: post.first_name + " " + post.last_name
      });
      setShowComments(true);
      // Fetch comments
      const res = await fetch(`http://localhost:8080/group/Getcomments/${post.id}/${offsetcomment.current}/${grpID}`, {
        method: "GET",
        credentials: "include",
      });

      if (!res.ok) {
        return false
      }
      const data = await res.json() || [];



      if (data.length == 0) {
        return false
      } else {
        offsetcomment.current += 10

      }


      setComment([...comment, ...data]);
      return data[0].id

    } catch (err) {
      return false
    }
    finally {
      setLoadingComment(false);
    }
  }

  return (
    <div>
      <PostCreationTrigger setPost={setPost} groupId={grpID} />
      {toast && (
        <div className={`toast ${toast.type}`}>
          <span>{toast.message}</span>
          <button onClick={() => setToast(null)} className="toast-close">×</button>
        </div>
      )}

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
              onGetComments={GetComments}
              ondolike={handleLike}
            />
          ))}

          <Comment
            comments={comment}
            isOpen={showComments}
            onClose={closeComments}
            postId={selectedPost?.id}
            postTitle={selectedPost?.title}
            onCommentChange={refreshComments}
            lodinggg={loadingcomment}
            ongetcomment={GetComments}
            post={selectedPost}
            showToast={showToast}
            ID={grpID}
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