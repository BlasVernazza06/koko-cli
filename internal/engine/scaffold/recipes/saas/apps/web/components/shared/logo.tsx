import { motion } from "motion/react"
import { Boxes } from "lucide-react"

export default function Logo() {
    return (
        <motion.div className="flex justify-center gap-2 md:justify-start" initial={{ x: -250 }} animate={{ x: 0 }} transition={{ duration: 1.2, ease: "easeInOut" }}>
          <a href="#" className="flex items-center gap-2 font-bold text-xl">
            <div className="flex size-10 items-center justify-center rounded-md bg-zinc-900 border border-zinc-880 text-zinc-400">
              <Boxes size={20} />
            </div>
            [[.ProjectName]]
          </a>
        </motion.div>
    );
}