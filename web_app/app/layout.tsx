import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { SSEStatusProvider } from "@services/contextService"; // Import SSEStatusProvider for managing SSE status updates across the app

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Manhattan Tech Ventures",
  description: "Created By Sarthak Joshi",
};

export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <SSEStatusProvider> {/* Wraps the application in SSEStatusProvider to provide SSE status updates */}
          {children} {/* Render the children components, which represent the rest of the application */}
        </SSEStatusProvider>
      </body>
    </html>
  );
}
