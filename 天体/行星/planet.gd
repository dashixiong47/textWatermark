extends Node2D
var calculate_size = preload("res://srcipt/calculate_size.gd").new()

@onready var sprite_2D= $Sprite2D
var roche_limit=100

func _ready() -> void:
	calculate_size.set_size(self,sprite_2D,roche_limit)
func _on_panel_mouse_entered() -> void:
	print(name)	
	pass # Replace with function body.
func  _on_panel_mouse_exited() -> void:
	print(name+"移出")
