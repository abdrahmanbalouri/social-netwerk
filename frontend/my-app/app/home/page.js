"use client";

import { useEffect, useState } from "react";
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

      <form onSubmit={logout}>
        <button type="submit">Logout</button>
      </form>
    </div>
  );
}
