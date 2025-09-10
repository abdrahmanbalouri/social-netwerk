"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";


export default function Profile() {
//  const [loading, setLoading] = useState(true);
  const router = useRouter();


  async function logout(e) {
    e.preventDefault();
    const res = await fetch("http://localhost:8080/api/logout", {
      method: "POST",
      credentials: "include",
    });
    if (!res.ok) {
     // alert(await res.text());
      return;
    }
    router.replace("/login"); 
  }

 // if (loading) return <p>Loading...</p>;
  //if (!user) return null;

  return (
    <div>
     
      <form onSubmit={logout}>
        <button type="submit">Logout</button>
      </form>
    </div>
  );
}
