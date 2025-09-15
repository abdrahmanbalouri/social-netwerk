"use client";
import './Home.css';
import { useRouter } from "next/navigation";



export default function home() {
  const router = useRouter();


  async function logout(e) {
    e.preventDefault();
    const res = await fetch("http://localhost:8080/api/logout", {
      method: "POST",
      credentials: "include",
    });
    console.log(res);
    
    if (!res.ok) {
      return;
    }
    router.replace("/login"); 
  }


return (
  <div>
    <nav className="navbar">
      <div class="Container">
        <h2 class="log">social</h2>
      </div>
      <div class="search-bar">
        <i class="fa-solid fa-magnifying-glass"></i>
        <input type="search" placeholder="search for some one"></input>
      </div>
      <div class="create">
          <label class="btn btn-primary" for="create-post">
            <i class="fa-solid fa-plus"></i>
          </label>
          <div class="profile-picture">
              <img src="/avatar.png" alt="Profile" />
          </div>
      </div>
    </nav>

    <main className="content">
     
    </main>
  </div>
);

}
