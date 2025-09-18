"use client";
import Providers from "./providers";
import "../app/home/Home.css"; // ...existing global imports if needed

export const metadata = {
  title: "Social Network",
  description: "App",
};

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>
        <Providers>
          {children}
        </Providers>
      </body>
    </html>
  );
}
