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
import { useWS } from "../../../context/wsContext.js";
import AddReactionIcon from "@mui/icons-material/AddReaction";
import { FileText, SendIcon } from "lucide-react";
import "../../../styles/chat.css";

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

      const data = await response.json() || [];

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

export function GroupChat({ groupId }) {
  console.log("group chat component rendered");
  const [id, setId] = useState(null);
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState("");
  const [showEmojis, setShowEmojis] = useState(false);
  const inputRef = useRef(null);
  const chatEndRef = useRef(null);
  const { sendMessage, addListener, removeListener } = useWS();
  useEffect(() => {
    fetch("http://localhost:8080/api/me", {
      credentials: "include",
      method: "GET",
    })
      .then((res) => res.json())
      .then((data) => {
        console.log("user id in group chat------------ : ", data.user_id);
        setId(data.user_id)
      })
      .catch((err) => console.error(err));
  });

  setTimeout(() => {
    inputRef.current?.focus();
  }, 0);

  useEffect(() => {
    const fetchMessages = async () => {
      try {
        const response = await fetch(
          `http://localhost:8080/api/getGroupMessages?groupId=${groupId}`,
          { credentials: "include", method: "GET" }
        );
        if (!response.ok) throw new Error("Failed to fetch messages");
        const data = await response.json();
        if (data.messages) {
          console.log("user id in group chat message listener23131 : ", id);

          const formattedMessages = data.messages
            .map((msg) => ({
              text: msg.content,
              sender: msg.senderId !== id ? "them" : "me",
              time: new Date(msg.createdAt).toLocaleString(),
            }))
            .reverse();

          setMessages(formattedMessages);
        }
      } catch (error) {
        console.error("Error fetching messages:", error);
      }
    };
    fetchMessages();
  }, [id]);

  useEffect(() => {
    const handleIncomingMessage = (data) => {
      console.log("user id in group chat message listener : ", id);
      setMessages((prev) => [
        ...prev,
        {
          text: data.content,
          sender: data.senderId !== id ? "them" : "me",
          time: data.time,
        },
      ]);
    };

    addListener("group_message", handleIncomingMessage);
    return () => removeListener("group_message", handleIncomingMessage);
  }, [addListener, removeListener, id]);

  // 👇 Scroll to bottom whenever messages change
  useEffect(() => {
    if (chatEndRef.current) {
      chatEndRef.current.scrollIntoView({ behavior: "smooth" });
    }
  }, [messages]);

  const emojiArray = [
    "😀",
    "😃",
    "😄",
    "😁",
    "😆",
    "😅",
    "🤣",
    "😂",
    "🚀",
    "💡",
    "😊",
    "😇",
    "🙂",
    "🙃",
    "😉",
    "😍",
    "🥰",
    "😘",
    "😗",
    "😋",
    "😛",
    "😜",
    "🤪",
    "😝",
    "🤑",
    "🤗",
    "🤭",
    "🤔",
    "🤨",
    "😐",
    "😑",
    "😶",
    "😏",
    "😒",
    "🙄",
    "😬",
    "😔",
    "😪",
    "🤤",
    "😴",
    "😷",
    "🤒",
    "🤕",
    "🤢",
    "🤮",
    "🥴",
    "😵",
    "🤯",
    "😎",
    "🤓",
  ];

  function handleSendGroupChatMessage() {
    console.log("=========", typeof groupId);
    if (input.trim() === "") return;
    const payload = {
      groupID: groupId,
      messageContent: input,
      type: "group_message",
    };
    sendMessage(payload);
    setInput("");
    setShowEmojis(false);
  }

  const addEmoji = (emoji) => {
    const cursorPos = inputRef.current.selectionStart;
    const newText = input.slice(0, cursorPos) + input.slice(cursorPos) + emoji;
    setInput(newText);
    setTimeout(() => {
      inputRef.current.focus();
      const end = inputRef.current.value.length;
      inputRef.current.setSelectionRange(end, end);
    }, 0);
  };

  return (
    <div className="group-chat-container">
      {/* Chat box */}
      <div className="chat-box">
        {messages.length === 0 ? (
          <p className="no-msg">No messages yet</p>
        ) : (
          messages.map((msg, index) => (
            <div key={index} className={`message ${msg.sender}`}>
              <div className="msg-content">{msg.text}</div>
              <div className="info-time">
                <span className="time">
                  {new Date(msg.time).toLocaleTimeString("en-US")}
                </span>
              </div>
            </div>
          ))
        )}
        <div ref={chatEndRef}></div>
      </div>

      {/* Emoji panel */}
      {showEmojis && (
        <div className="emoji-panel">
          {emojiArray.map((emoji, i) => (
            <span key={i} className="emoji" onClick={() => addEmoji(emoji)}>
              {emoji}
            </span>
          ))}
        </div>
      )}

      {/* Input area */}
      <div className="input-area">
        <button
          className="emoji-toggle"
          onClick={() => setShowEmojis(!showEmojis)}
        >
          <AddReactionIcon />
        </button>

        <input
          ref={inputRef}
          type="text"
          placeholder="Type a message..."
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && handleSendGroupChatMessage()}
        />
        <button onClick={handleSendGroupChatMessage}>
          <SendIcon />
        </button>
      </div>
    </div>
  );
}