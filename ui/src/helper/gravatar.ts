export async function gravatarUrl(identifier: string, size = 40): Promise<string> {
  const clean = identifier.trim().toLowerCase();
  const buffer = await crypto.subtle.digest('SHA-256', new TextEncoder().encode(clean));
  const hash = Array.from(new Uint8Array(buffer)).map(b => b.toString(16).padStart(2, '0')).join('');
  return `https://gravatar.com/avatar/${hash}?d=mp&s=${size}`;
}
