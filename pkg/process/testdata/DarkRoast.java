package com.javaprogramto.programs.infinite.loops;

public class DarkRoast {
	public static void main(String[] args) {
		try {
			while (true) {
				Thread.sleep(4000);
				System.out.println("Running while loop");
			}
		} catch (InterruptedException e) {
			System.out.println("Interrupted");
		}
	}
}
