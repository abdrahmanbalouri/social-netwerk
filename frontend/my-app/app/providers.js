"use client";
import { DarkModeProvider } from "../context/darkMod";
import { ProfileProvider } from "../context/profile";

export default function Providers({ children }) {
  return (
    <DarkModeProvider>
      <ProfileProvider>{children}</ProfileProvider>
    </DarkModeProvider>
  );
}
