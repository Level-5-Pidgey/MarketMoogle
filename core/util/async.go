/*
 * Copyright (c) 2022 Carl Alexander Bird.
 * This file (async.go) is part of MarketMoogle and is released GNU General Public License.
 * Please see the "LICENSE" file within MarketMoogle to view the full license. This file, and all code within MarketMoogle fall under the GNU General Public License.
 */

package util

func Async[T any](f func() T) chan T {
	ch := make(chan T)

	go func() {
		ch <- f()
	}()

	return ch
}
