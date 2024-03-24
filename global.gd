extends Node

# 引力常数
const G = 6.67430e-3*10000

#func _physics_process(delta):
	#var bodies = get_tree().get_nodes_in_group("bodies")
	#
	#for body1 in bodies:
		#for body2 in bodies:
			#if body1 == body2:
				#continue
			#
			#var distance_vector = body1.position - body2.position
			#var distance = distance_vector.length()
			#if distance == 0:
				#continue
			#
			#var force_magnitude = G * body1.mass * body2.mass / (distance * distance)
			#var force_direction = distance_vector.normalized()
			#var roche_limit = 2.5 * (body1.roche_limit + body2.roche_limit)
			## 旋转前的简单模拟
			#if distance > roche_limit:
				#var tangent_direction = Vector2(-force_direction.y, force_direction.x)
				#body1.linear_velocity += tangent_direction * sqrt(force_magnitude / body1.mass)
				#body2.linear_velocity += -tangent_direction * sqrt(force_magnitude / body2.mass)
			#else:
				## 达到洛希极限后的行为，如简单地让它们直接碰撞或其他交互
				#body1.apply_central_force(-force_direction * force_magnitude)
				#body2.apply_central_force(force_direction * force_magnitude)
