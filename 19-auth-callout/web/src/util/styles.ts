export default function cn(...args: any[]): string {
  return args.filter(Boolean).join(' ');
}
