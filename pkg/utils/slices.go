/*
 * Copyright (C) 2019  Appvia Ltd <info@appvia.io>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package utils

// Contains checks a list has a value in it
func Contains(v string, l []string) bool {
	for _, x := range l {
		if v == x {
			return true
		}
	}

	return false
}

// RemoveItem removes a value from a slice
func RemoveItem(v string, list []string) []string {
	var a []string

	for _, x := range list {
		if x != v {
			a = append(a, x)
		}
	}

	return a
}

// ChunkBy breaks a slice into chunks
func ChunkBy(items []string, chunkSize int) (chunks [][]string) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}

	return append(chunks, items)
}
