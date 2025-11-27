import './globals.css';

import type { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'Veritas - Real-Time News & Search  Agent',
  description: 'Answer user questions with accurate, real-time information',
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
