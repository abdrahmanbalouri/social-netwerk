"use client";
import React, { useEffect, useRef, useState } from "react";
import "../../styles/gallery.css";
import { useProfile } from "../../context/profile";
import { useDarkMode } from "../../context/darkMod";
import Navbar from "../../components/Navbar";
import LeftBar from "../../components/LeftBar";
import RightBar from "../../components/RightBar";
import { middleware } from "../../middleware/middelware";
import { useWS } from "../../context/wsContext";
import { useRouter } from "next/navigation";

export default function Gallery() {
  const [images, setImages] = useState([]);
  const [imgIndex, setImgIndex] = useState("");
  const slideRef = useRef(null);
  const { Profile } = useProfile();
  const { darkMode } = useDarkMode();
  const router = useRouter();
  const sendMessage = useWS()
  // Authentication check
  useEffect(() => {
    const checkAuth = async () => {
      const auth = await middleware();
      if (!auth) {
        router.push("/login");
        sendMessage({ type: "logout" })
      }
    }
    checkAuth();
  }, [])

  useEffect(() => {
    if (!Profile?.id) return;

    fetch(`http://localhost:8080/api/gallery?id=${Profile.id}`, {
      credentials: "include",
    })
      .then((res) => res.json())
      .then((data) => {

        if (data) {
          let images = data.filter(img => img.imagePath);
          setImages(images);
          return
        }
        if (!data ||  data.length > 0) {
          //setImgIndex(`url(/${data[0].imagePath})`);
        }
      })
      .catch((err) => console.error(err));
  }, [Profile]);

  function next() {
    let items = slideRef.current.querySelectorAll(".gallery-item");
    if (items.length === 0) return;

    slideRef.current.appendChild(items[0]);

    const bg = window.getComputedStyle(items[items.length - 1]).backgroundImage;
    setImgIndex(bg);
  }

  function prev() {
    let items = slideRef.current.querySelectorAll(".gallery-item");
    if (items.length === 0) return;

    slideRef.current.prepend(items[items.length - 1]);

    const bg = window.getComputedStyle(items[items.length - 1]).backgroundImage;
    setImgIndex(bg);
  }

  return (
    <div className={darkMode ? "theme-dark" : "theme-light"}>
      <Navbar />
      <div className="gallery-page">
        <LeftBar showSidebar={true} />
        <main className="gallery-main">
          <div className="gallery-wrapper" style={{ backgroundImage: imgIndex }}>
            <div className="gallery-slider" ref={slideRef}>
              {images.length === 0 ? (
                <div
                  className="gallery-item"
                  style={{ backgroundImage: `url(/assets/default.png)` }}
                ></div>
              ) : (
                images.map((img, index) => (
                  <div
                    className="gallery-item"
                    key={index}
                    style={{ backgroundImage: `url(/${img.imagePath})` }}
                  >
                    <div className="gallery-content">
                      <div className="gallery-title">{img.title || "No Title"}</div>
                      <div className="gallery-desc">
                        {img.description || "Lorem ipsum dolor sit amet consectetur adipisicing elit. Ab, eum!"}
                      </div>
                    </div>
                  </div>
                ))
              )}

            </div>

            <div className="gallery-nav">
              <button className="gallery-prev" onClick={prev}>
                <i className="fa-solid fa-arrow-left"></i>
              </button>
              <button className="gallery-next" onClick={next}>
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
