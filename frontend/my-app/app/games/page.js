"use client";
import React, { useEffect } from "react";
import Navbar from "../../components/Navbar.js";
import LeftBar from "../../components/LeftBar.js";
import RightBar from "../../components/RightBar.js";
import { useDarkMode } from "../../context/darkMod.js";
import Link from "next/link.js";
import "../../styles/games.css";
import { useWS } from "../../context/wsContext.js";

export default function Game() {
  const sendMessage = useWS()
  // Authentication check
  useEffect(() => {
    const checkAuth = async () => {
      const midle = await middleware();
      return midle;
    }
    const auth = checkAuth()
    if (!auth) {
      router.push("/login");
      sendMessage({ type: "logout" })
    }
  }, [])
  const { darkMode } = useDarkMode();
  const [showSidebar, setShowSidebar] = React.useState(true);
  return (
    <div className={darkMode ? "theme-dark" : "theme-light"}>
      <Navbar
        onToggleSidebar={() => setShowSidebar(!showSidebar)}
      />
      <main className="content">
        <LeftBar showSidebar={showSidebar} />

        <div className="game-container">
          <h2 className="game-title">Mini Games</h2>
          <p className="game-desc">Enjoy a selection of simple games below!</p>
          <div className="games-list">
            <div className="game-card">
              <h3>Boomber man</h3>
              <Link href="https://raji383.github.io/bomberman/" target="_blank">
                <img src="/bomberman.png" alt="Boomber man" />
              </Link>
            </div>
            <div className="game-card">
              <h3>Arkanoid</h3>
              <Link href="https://brick-breaker-abaid.netlify.app/" target="_blank">
                <img src="/Arkanoid.png" alt="Arkanoid" />
              </Link>
            </div>
            <div className="game-card">
              <h3>super-mario</h3>
              <Link href="https://raji383.github.io/super-mario/" target="_blank">
                <img src="/super-mario.png" alt="Arkanoid" />
              </Link>
            </div>
            <div className="game-card">
              <h3>Raven</h3>
              <Link href="https://raji383.github.io/Raven-game/" target="_blank">
                <img src="/Raven.png" alt="Arkanoid" />
              </Link>
            </div>
            <div className="game-card">
              <h3>shadow</h3>
              <Link href="https://raji383.github.io/shadow-game/" target="_blank">
                <img src="/shadow.png" alt="Arkanoid" />
              </Link>
            </div>
            <div className="game-card">
              <h3>War-machine</h3>
              <Link href="https://raji383.github.io/War-machine/" target="_blank">
                <img src="/war.png" alt="Arkanoid" />
              </Link>
            </div>
            <div className="game-card">
              <h3>toad</h3>
              <Link href="https://raji383.github.io/War-machine/" target="_blank">
                <img src="/toad.png" alt="Arkanoid" />
              </Link>
            </div>
          </div>
        </div>

        <RightBar />
      </main>
    </div>
  );
}


