import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

/**
 * Get the full avatar URL
 * If avatar is already a full URL (starts with http:// or https://), return as is
 * Otherwise, prepend the API base URL
 */
export function getAvatarUrl(avatar: string | null | undefined): string | undefined {
  if (!avatar) return undefined;
  
  // Check if avatar is already a full URL (GitHub, external URL, etc.)
  if (avatar.startsWith('http://') || avatar.startsWith('https://')) {
    return avatar;
  }
  
  // Otherwise, it's a local upload - prepend API URL
  return `${import.meta.env.VITE_API_URL}${avatar}`;
}
