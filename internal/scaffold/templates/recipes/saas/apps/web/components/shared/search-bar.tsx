import { Input } from "@/components/ui/input";
import { Search } from "lucide-react";

interface SearchBarProps {
    value: string;
    onChange: (value: string) => void;
    placeholder?: string;
}

export default function SearchBar({ value, onChange, placeholder = "Buscar..." }: SearchBarProps) {
    return (
        <div className="relative flex-1">
            <Search className="absolute left-4 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input 
                placeholder={placeholder}
                className="pl-10 h-12 bg-white dark:bg-zinc-900 border-border/40 rounded-2xl focus:ring-brand/10 transition-all font-medium text-sm"
                value={value}
                onChange={(e) => onChange(e.target.value)}
            />
        </div>
    );
}