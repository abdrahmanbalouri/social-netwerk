"use client";
import Navbar from "../../../components/Navbar.js"
import { GroupPostChat } from "../../../components/groupPostCat.js";
import Post from "../../../components/Post.js";
import { useEffect, useState } from "react";
import "../../../styles/groupstyle.css"
import { useParams } from "next/navigation";
import { PostCreationTrigger } from "../../../components/cretaePostGroup.js"
import LeftBar from "../../../components/LeftBar.js";
import RightBar from "../../../components/RightBar.js";
// import {CreatePost} from ""


export default function () {
    return (
        <>
            <Navbar />
            <main className="content">
                <LeftBar showSidebar={true} />
                {/* <AllPosts /> */}
                <GroupPostChat />
                <RightBar />
            </main>
            {/* <AllPosts /> */}
            {/* <CreatePost /> */}
        </>
    )
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
            method: 'GET',
            credentials: 'include',
        })
            .then(res => {
                if (!res.ok) throw new Error('Failed to fetch posts');
                return res.json();
            })
            .then(data => {
                setPost(data);
                setLoading(false);
            })
            .catch(error => {
                console.error("Failed to fetch posts:", error);
                setLoading(false);
            });
    }, [grpID]);
    if (!posts) {
        return (
            <>
                <PostCreationTrigger />
                <div>There is no post yeeeeeet.</div>
            </>
        );
    }

    console.log("posts are :", posts);
    return (
        <div>
            <PostCreationTrigger />
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
        </div>
    )
}

function GetComments(post) {
    // setSelectedPost({
    //     id: post.id,
    //     title: post.title || post.post_title || "Post"
    // });
}

function AddLike() {

}

export async function CreatePost(groupId, formData) {
    console.log("grouuup id is :", groupId);

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
