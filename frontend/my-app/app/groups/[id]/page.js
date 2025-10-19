"use client";
import Navbar from "../../../components/Navbar.js"
import Post from "../../../components/Post.js";
import { useEffect, useState } from "react";
import "./page.css"
import { useParams } from "next/navigation";
// import {CreatePost} from ""


export default function () {
    return (
        <>
            <Navbar />
            <AllPosts />
        </>
    )
}

function AllPosts() {

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
        return <div>There is no post yetttttt.</div>;
    }
    console.log("posts are :", posts);
    return (
        <>
            {/* <CreatePost /> */}
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
        </>
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
