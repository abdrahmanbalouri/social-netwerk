"use client";
import { createContext, useContext, useEffect, useState } from "react";

const ProfileContext = createContext(undefined);

export function ProfileProvider({ children }) {
  const [profile, setProfile] = useState(null);

  async function loadProfile() {
    try {
      const res = await fetch("http://localhost:8080/api/profile?userId=0", {
        method: "GET",
        credentials: "include",
      });
      if (res.ok) {
        const json = await res.json();
        setProfile(json);
      }
    } catch (err) {
      console.error("loadProfile", err);
    }
  }

  useEffect(() => {
    loadProfile();
  }, []);

  return (
    <ProfileContext.Provider value={{ profile, setProfile, reload: loadProfile }}>
      {children}
    </ProfileContext.Provider>
  );
}

export function useProfile() {
  const ctx = useContext(ProfileContext);
  if (ctx === undefined) {
    throw new Error("useProfile must be used within a ProfileProvider");
  }
  return ctx;
}
