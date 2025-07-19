import {  clsx } from "clsx"
import { twMerge } from "tailwind-merge"
import type {ClassValue} from "clsx";

export function cn(...inputs: Array<ClassValue>) {
  return twMerge(clsx(inputs))
}

export function formatToDDMMYY(dateString: string): string {
  const date = new Date(dateString);
  const day = String(date.getUTCDate()).padStart(2, '0');
  const month = String(date.getUTCMonth() + 1).padStart(2, '0');
  const year = String(date.getUTCFullYear()).slice(-2);
  return `${day}-${month}-${year}`;
}