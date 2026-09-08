import mongoose, { Schema, Document } from 'mongoose';

export interface IItem extends Document {
  title: string;
  description?: string;
  status: 'pending' | 'in_progress' | 'completed';
  createdAt: Date;
  updatedAt: Date;
}

const ItemSchema: Schema = new Schema(
  {
    title: { type: String, required: true, trim: true },
    description: { type: String, default: '' },
    status: {
      type: String,
      enum: ['pending', 'in_progress', 'completed'],
      default: 'pending',
    },
  },
  { timestamps: true }
);

export const Item = mongoose.models.Item || mongoose.model<IItem>('Item', ItemSchema);
