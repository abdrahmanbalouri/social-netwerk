"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";



export default function Profile() {
  const router = useRouter();

  useEffect(() => {
    async function fetchUser() {
      try {
        const response = await fetch("http://localhost:8080/api/me", {
          credentials: "include",
        });

        if (!response.ok) {
          router.replace("/login"); // redirect if not authenticated
          return null;
        } else if (response.ok) {

          router.replace("/home");
        }
      } catch (error) {
        router.replace("/login");
        return null;

      }
    }
    fetchUser()

  });

  async function logout(e) {
    e.preventDefault();
    const res = await fetch("http://localhost:8080/api/logout", {
      method: "POST",
      credentials: "include",
    });
    if (!res.ok) {
      return;
    }
    router.replace("/login"); // redirect to login page
  }


  return (
    <div>

      <form onSubmit={logout}>
        <button type="submit">Logout</button>
      </form>
    </div>
  );
}
