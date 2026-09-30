import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';
import { IsNotEmpty, IsOptional, IsString, IsIn } from 'class-validator';

export class CreateTaskDto {
  @ApiProperty({ description: 'Title of the task', example: 'Setup CI/CD pipeline' })
  @IsString()
  @IsNotEmpty()
  title: string;

  @ApiPropertyOptional({ description: 'Detailed description', example: 'Automate deployment' })
  @IsString()
  @IsOptional()
  description?: string;

  @ApiPropertyOptional({ enum: ['pending', 'in_progress', 'completed'], default: 'pending' })
  @IsIn(['pending', 'in_progress', 'completed'])
  @IsOptional()
  status?: string;

  @ApiPropertyOptional({ enum: ['low', 'medium', 'high', 'urgent'], default: 'medium' })
  @IsIn(['low', 'medium', 'high', 'urgent'])
  @IsOptional()
  priority?: string;
}

export class UpdateTaskDto {
  @ApiPropertyOptional()
  @IsString()
  @IsOptional()
  title?: string;

  @ApiPropertyOptional()
  @IsString()
  @IsOptional()
  description?: string;

  @ApiPropertyOptional({ enum: ['pending', 'in_progress', 'completed'] })
  @IsIn(['pending', 'in_progress', 'completed'])
  @IsOptional()
  status?: string;

  @ApiPropertyOptional({ enum: ['low', 'medium', 'high', 'urgent'] })
  @IsIn(['low', 'medium', 'high', 'urgent'])
  @IsOptional()
  priority?: string;
}
