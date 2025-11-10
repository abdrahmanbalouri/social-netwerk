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
import EventCard from "../../../components/EventCard.js";
import { useWS } from "../../../context/wsContext.js";
import AddReactionIcon from "@mui/icons-material/AddReaction";
import SendIcon from "@mui/icons-material/Send";
import "../../../styles/chat.css";
import { useProfile } from "../../../context/profile.js";
import ShowToast from "../../../components/ShowToast.js";


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
        InvitationType: "invitation",
        invitedUsers: [invitedUserID],
      }),
    });

    const temp = await res.json();
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

          <EventForm ShowEventForm={ShowEventForm} closeForm={showEvent}
            fetchEvents={fetchEvents} showEvent={showEvent}
          />
        </div>
      )}

      {EventList?.length === 0 ? (
        <div className="no-data">No events yet. Be the first to create one!</div>
      ) : (
        EventList.map((ev) => (
          <EventCard key={ev.id} ev={ev} goingEvent={goingEvent} />
        ))
      )}
    </>
  )
}


export function EventForm({ closeForm, fetchEvents, showEvent }) {
  const params = useParams();
  const grpID = params.id;
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [dateTime, setDateTime] = useState("");
  const [error, setError] = useState("");
  const { sendMessage } = useWS();


  async function createEvent(e) {
    e.preventDefault()
    const Data = {
      title: title,
      description: description,
      dateTime: dateTime
    }

    if (!title || !description || !dateTime) {
      setError("All fields are required");
      return;
    }

    if ((title.length < 5 || title.length > 50)) {
      setError(" Title must be between 5 and 50 characters");
      return;
    }

    if (description.length < 10 || description.length > 300) {
      setError("Description must be between 10 and 300 characters");
      return;
    }


    if (new Date(dateTime) <= new Date()) {
      setError("Event date and time must be in the future");
      return;
    }

    try {
      const res = await fetch(`http://localhost:8080/api/createEvent/${grpID}`, {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(Data),
      });

      if (!res.ok) {
        let err = await res.text();
        setError(err);
        return;
      }

      await fetchEvents();

      closeForm();


      sendMessage({
        type: "new event",
        receiverId: params.id,
        messageContent: "",
      });


    } catch (error) {
      setError(error.message);
    }
  }



  return (

    <div className="create-event-card">
      <h2 className="card-title">Create Event</h2>
      <form className="event-form">
        <div className="form-grp">
          <label htmlFor="event-title">Title</label>
          <input type="text" id="event-title" placeholder="Event title..." onChange={(e) => { setTitle(e.target.value) }} />
        </div>

        <div className="form-grp">
          <label htmlFor="event-description">Description</label>
          <textarea id="event-description" rows="3" placeholder="Event description..." onChange={(e) => { setDescription(e.target.value) }}></textarea>
        </div>

        <div className="form-grp">
          <label htmlFor="event-datetime">Day/Time</label>
          <input type="datetime-local" id="event-datetime" onChange={(e) => { setDateTime(e.target.value) }} />
        </div>
        <ShowToast key={Date.now()} message={error} />
        <button type="submit" className="btn-create" onClick={createEvent}>Create Event</button>
        <button type="button" className="cancel-btn" onClick={showEvent}>Cancel</button>
      </form>
    </div>


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
  const [error, seterror] = useState(null);
  const offsetcomment = useRef(0);
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

        if (data?.error) {
          seterror(data.error)
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
      if (data?.error) {
        seterror(data.error)
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
      const res = await fetch(`http://localhost:8080/group/like/${postId}/${grpID}`, {
        method: "POST",
        credentials: "include",
      });
      const response = await res.json();


      if (response.error) {
        if (response.error === "Unauthorized") {
          router.push("/login");
        }
        return;
      }

      const updatedPost = await fetchPostById(postId);


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
      <ShowToast key={Date.now()} message={error} />

      {loading ? (
        <div>Loading posts...</div>
      ) : posts && posts.length === 0 ? (
        <div className="no-data">There is no post yet.</div>
      ) : (
        <div className="posts-list">
          {posts && posts.map((post) => (
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

    const text = await res.json();

    if (text.error) {
      return text
    }

    return text;
  } catch (error) {
    console.error("CreatePost error:", error);
  }
}

export function GroupChat({ groupId }) {
  const [id, setId] = useState(null);
  const [preview, setPreview] = useState(null);
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState("");
  const [showEmojis, setShowEmojis] = useState(false);
  const [Error, setError] = useState(null)
  const inputRef = useRef(null);
  const chatEndRef = useRef(null);
  const { sendMessage, addListener, removeListener } = useWS();
  const { Profile } = useProfile();

  const scrollToBottom = () => {
    if (chatEndRef.current) {
      chatEndRef.current.scrollIntoView({ behavior: "smooth" });
    }
  };
  useEffect(() => {
    fetch("http://localhost:8080/api/me", {
      credentials: "include",
      method: "GET",
    })
      .then((res) => res.json())
      .then((data) => {

        setId(data.user_id)
      })
      .catch((err) => setError(err));
  });

  setTimeout(() => {
    inputRef.current?.focus();
  }, 0);

  const handleFileChange = (e) => {
    const file = e.target.files[0];
    if (!file) {
      setPreview(null);
      return;
    }

    const allowedTypes = ["image/png", "image/jpeg", "image/jpg", "image/gif", "image/webp"];
    if (!allowedTypes.includes(file.type)) {
      setError("Only image files are allowed!");
      e.target.value = "";
      setPreview(null);
      return;
    }

    const reader = new FileReader();
    reader.onloadend = () => {
      setPreview(reader.result);
    };
    reader.readAsDataURL(file);
  };



  const removeImage = () => {
    setPreview(null);
    document.getElementById("idfile").value = "";
  };

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

          const formattedMessages = data.messages
            .map((msg) => ({
              text: msg.content,
              sender: msg.senderId !== id ? "them" : "me",
              time: new Date(msg.createdAt).toLocaleString(),
              name: msg.first_name + " " + msg.last_name,
              image: msg.photo,
              PictureSend: msg.PictureSend
            }))
            .reverse();

          setMessages(formattedMessages);
        }
      } catch (error) {
        setError(error);
      }
    };
    fetchMessages();
  }, [id]);

  useEffect(() => {
    const handleIncomingMessage = (data) => {

      setMessages((prev) => [
        ...prev,
        {
          text: data.content,
          sender: data.from !== id ? "them" : "me",
          time: data.time,
          name: data.name,
          image: data.image,
          PictureSend: data.PictureSend
        },
      ]);
    };

    addListener("group_message", handleIncomingMessage);
    return () => removeListener("group_message", handleIncomingMessage);
  }, [addListener, removeListener, id]);

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const emojiArray = ["😀", "😃", "😄", "😁", "😆", "😅", "🤣", "😂", "🚀", "💡", "😊", "😇", "🙂", "🙃", "😉", "😍", "🥰", "😘", "😗", "😋", "😛", "😜", "🤪", "😝", "🤑", "🤗", "🤭", "🤔", "🤨", "😐", "😑", "😶", "😏", "😒", "🙄", "😬", "😔", "😪", "🤤", "😴", "😷", "🤒", "🤕", "🤢", "🤮", "🥴", "😵", "🤯", "😎", "🤓"];

  const handleSendGroupChatMessage = () => {
    if (input.trim() === "" && !preview) return;
    if (input.length > 1000) {
      setError("message is too long")
      return
    }
    const payload = {
      groupID: groupId,
      messageContent: input,
      type: "group_message",
      name: Profile.first_name,
      image: Profile.image,
      PictureSend: preview
    };
    sendMessage(payload);
    setInput("");
    removeImage()
    setShowEmojis(false);
    setTimeout(() => scrollToBottom(), 100);
  };

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
<ShowToast key={Date.now()} message={Error}/>
      {/* Chat box */}
      <div className="chat-box">
        {messages.length === 0 ? (
          <p className="no-msg">No messages yet</p>
        ) : (
          messages.map((msg, index) => (
            <div key={index} className={`message ${msg.sender}`}>
              <div className="message-header">
                <img src={msg?.image ? `/uploads/${msg?.image}` : "/assets/default.png"} alt="profile picture" className="profilePic" />
                <span className="message-name">{msg.name}</span>
              </div>
              <div className="message-bubble">
                {msg.PictureSend && (
                  <img
                    src={`/uploads/${msg.PictureSend}`}
                    alt="message image"
                    className="message-img"
                  />
                )}
                <div className="msg-content">{msg.text}</div>
              </div>
              <div className="info-time">
                <span className="time">
                  {new Date(msg.time).toLocaleTimeString("en-US", {
                    hour: 'numeric',
                    minute: '2-digit'
                  })}
                </span>
              </div>
            </div>
          ))
        )}
        {preview && (
          <div className="preview-container">
            <img src={preview} alt="preview" className="preview-image" />
            <button className="remove-btn" onClick={removeImage}>
              ×
            </button>
          </div>
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
        <input
          hidden
          type="file"
          id="idfile"
          accept="image/*"
          onChange={handleFileChange}
        />
        <label htmlFor="idfile">
          <div className="emoji-toggle">
            <i className="fa-solid fa-image " ></i>
          </div>
        </label>
        <button onClick={handleSendGroupChatMessage}>
          <SendIcon />
        </button>
      </div>
    </div>
  );
}