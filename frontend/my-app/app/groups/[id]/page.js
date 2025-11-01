"use client";
import Navbar from "../../../components/Navbar.js";
import { GroupPostChat } from "../../../components/groupPostCat.js";
import Post from "../../../components/Post.js";
import { useEffect, useState } from "react";
import "../../../styles/groupstyle.css";
import { useParams } from "next/navigation";
import { PostCreationTrigger } from "../../../components/cretaePostGroup.js";
import LeftBar from "../../../components/LeftBar.js";
import RightBarGroup from "../../../components/RightBarGroups.js";
import { useDarkMode } from "../../../context/darkMod.js";
import { FileText, MessageCircle } from "lucide-react";
// import {CreatePost} from ""

export default function () {
  const { darkMode } = useDarkMode();

  return (
    <div className={darkMode ? "theme-dark" : "theme-light"}>
      <Navbar />
      <main className="content">
        <LeftBar showSidebar={true} />
        {/* <AllPosts /> */}
        <GroupPostChat />
        {/* <RightBar /> */}
        <RightBarGroup onClick={sendRequest} />
      </main>
      {/* <AllPosts /> */}
      {/* <CreatePost /> */}
    </div>
  );
}

function sendRequest(invitedUserID, grpID) {
  console.log("invited user id is : ", invitedUserID);
  console.log("group id is : ", grpID);

  fetch(`http://localhost:8080/group/invitation/${grpID}`, {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      invitedUsers: [invitedUserID],
    }),
  })
    .then(async (res) => {
      const temp = await res.json();
      console.log("temp is :", temp);
    })
    .catch((error) => {
      console.log("error sending invitation to the user :", error);
    });
}

export function AllPosts() {
  const [posts, setPost] = useState([]);
  const [loading, setLoading] = useState(true);
  const params = useParams();
  const grpID = params.id;

  useEffect(() => {
    if (!grpID) {
      setLoading(false);
      return;
    }

    fetch(`http://localhost:8080/group/fetchPosts/${grpID}`, {
      method: "GET",
      credentials: "include",
    })
      .then((res) => {
        if (!res.ok) throw new Error("Failed to fetch posts");
        return res.json();
      })
      .then((data) => {
        setPost(data || []);
        setLoading(false);
      })
      .catch((error) => {
        console.error("Failed to fetch posts:", error);
        setLoading(false);
      });
  }, []);
  // if (!posts) {
  //     return (
  //         <div>
  //             <PostCreationTrigger setPost = {setPost}/>
  //             <div>There is no post yeeeeeet.</div>
  //         </div>
  //     );
  // }

  console.log("posts are :", posts);
  return (
    <div>
      <PostCreationTrigger setPost={setPost} />
      <div className="content-area">
        {posts.length === 0 ? (
          <div className="group-empty-state">
            <FileText />
            <p className="group-empty-state-text">
              No posts yet. Be the first to share something!
            </p>
          </div>
        ) : (
          posts.map((post) => (
            <Post
              key={post.id}
              post={post}
              onGetComments={GetComments}
              ondolike={AddLike}
            />
          ))
        )}
      </div>
    </div>
  );
}
export function LastPost() {
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const params = useParams();
  const grpID = params.id;

  useEffect(() => {
    if (!grpID) {
      setLoading(false);
      return;
    }

    fetch(`http://localhost:8080/group/fetchPost/${grpID}`, {
      method: "GET",
      credentials: "include",
    })
      .then((res) => {
        if (!res.ok) throw new Error("Failed to fetch last post");
        return res.json();
      })
      .then((data) => {
        setPosts(data);
        setLoading(false);
      })
      .catch((error) => {
        console.error("Failed to fetch last post:", error);
        setLoading(false);
      });
  }, [grpID]);

  if (loading) {
    return <div>Loading...</div>;
  }

  if (!posts || posts.length === 0) {
    return (
      <div>
        <PostCreationTrigger setPosts={setPosts} />
        <div>There is no post yet.</div>
      </div>
    );
  }

  console.log("posts are:", posts);

  return (
    <div>
      <PostCreationTrigger setPosts={setPosts} />
      <div className="content-area">
        {posts.map((post) => (
          <Post
            key={post.id}
            post={post}
            onGetComments={GetComments}
            ondolike={AddLike}
          />
        ))}
      </div>
    </div>
  );
}

function GetComments(post) {
  // setSelectedPost({
  //     id: post.id,
  //     title: post.title || post.post_title || "Post"
  // });
}

function AddLike() { }

export async function CreatePost(groupId, formData) {
  const data = new FormData();
  data.append("postData", JSON.stringify(formData));

  const res = await fetch(`http://localhost:8080/group/addPost/${groupId}`, {
    method: "POST",
    credentials: "include",
    body: data,
  });

  const text = await res.text();

  if (!res.ok) throw new Error("Failed to create a post for groups");

  return JSON.parse(text);
}

// group chat component
export function GroupChat() {
  console.log("group chat component rendered");
  const [id, setId] = useState(null)
  const [input, setInput] = useState("");
  const [messages, setMessages] = useState([]);
  const { addListener, removeListener, sendMessage } = useWS();
  const emojiArray = ["😀", "😃", "😄", "😁", "😆", "😅", "🤣", "😂", "🚀", "💡", "😊", "😇", "🙂", "🙃", "😉", "😍", "🥰", "😘", "😗", "😋", "😛", "😜", "🤪", "😝", "🤑", "🤗", "🤭", "🤔", "🤨", "😐", "😑", "😶", "😏", "😒", "🙄", "😬", "😔", "😪", "🤤", "😴", "😷", "🤒", "🤕", "🤢", "🤮", "🥴", "😵", "🤯", "😎", "🤓"];
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
  // 👇 Scroll to bottom whenever messages change
  useEffect(() => {
    if (chatEndRef.current) {
      chatEndRef.current.scrollIntoView({ behavior: "smooth" });
    }
  }, [messages]);

  useEffect(() => {
    const handleIncomingMessage = (data) => {

      setMessages((prev) => [
        ...prev,
        {
          text: data.content,
          sender: data.from === user.id ? "them" : "me",
          time: data.time,
        },
      ]);

    };

    addListener("group_messages", handleIncomingMessage);
    return () => removeListener("group_messages", handleIncomingMessage);
  }, [addListener, removeListener]);
  useEffect(() => {
    fetch("http://localhost:8080/api/me")
      .then(res => res.json())
      .then(data => setId(data.user_id))
      .catch(err => console.error(err))
  })
  
  function handleSendGroupChatMessage() {
    if (input.trim() === "") return
    const payload = {
      senderId: id,
      messageContent: input,
      type: "message",
    }
    sendMessage(payload)


  }

  return (
    <div className="chat-container">
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
