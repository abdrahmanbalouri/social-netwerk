"use client";
import React, { useEffect, useRef, useState } from "react";
import "./gallery.css";
import { useProfile } from "../../context/profile";
import { useDarkMode } from "../../context/darkMod";
import Navbar from "../../components/Navbar";
import LeftBar from "../../components/LeftBar";
import RightBar from "../../components/RightBar";

export default function Gallery() {
  const [images, setImages] = useState([]);
  const [imgIndex, setImgIndex] = useState("");
  const slideRef = useRef(null);
  const { Profile } = useProfile();
  const { darkMode } = useDarkMode();

  useEffect(() => {
    if (!Profile?.id) return;

    fetch(`http://localhost:8080/api/gallery?id=${Profile.id}`, {
      credentials: "include",
    })
      .then((res) => res.json())
      .then((data) => {
        setImages(data);

        if (data.length > 0) {
          setImgIndex(`url(/${data[0].imagePath})`);
        }
      })
      .catch((err) => console.error(err));
  }, [Profile]);

  function next() {
    let items = slideRef.current.querySelectorAll(".item");
    slideRef.current.appendChild(items[0]);

    const bg = window.getComputedStyle(items[0]).backgroundImage;
    setImgIndex(bg); 
  }

  function prev() {
    let items = slideRef.current.querySelectorAll(".item");
    slideRef.current.prepend(items[items.length - 1]);

    const bg = window.getComputedStyle(items[items.length - 1]).backgroundImage;
    setImgIndex(bg); 
  }

  return (
    <div className={darkMode ? "theme-dark" : "theme-light"}>
      <Navbar />
      <div className="con" >
        <LeftBar showSidebar={true} />
        <main className="gallery">
          <div className="container" style={{ backgroundImage: imgIndex }}>
            <div className="slide" ref={slideRef}>
              {images.map((img, index) => (
                <div
                  className="item"
                  key={index}
                  style={{ backgroundImage: `url(/${img.imagePath})` }}
                >
                  <div className="content">
                    <div className="name">{img.title || "No Title"}</div>
                    <div className="des">
                      {img.description ||
                        "Lorem ipsum dolor sit amet consectetur adipisicing elit. Ab, eum!"}
                    </div>
                    <button>See More</button>
                  </div>
                </div>
              ))}
            </div>

            <div className="button">
              <button className="prev" onClick={prev}>
                <i className="fa-solid fa-arrow-left"></i>
              </button>
              <button className="next" onClick={next}>
                <i className="fa-solid fa-arrow-right"></i>
              </button>
            </div>
          </div>
        </main>
        <RightBar showSidebar={true} />
      </div>
    </div>
  );
}
